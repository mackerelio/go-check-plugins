package checklog

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/mackerelio/golib/pluginutil"
	"github.com/mattn/go-encoding"
	"github.com/mattn/go-zglob"
	"github.com/natefinch/atomic"
	enc "golang.org/x/text/encoding"
)

// overwritten with syscall.SIGTERM on unix environment (see check-log_unix.go)
var defaultSignal = os.Interrupt

type logOpts struct {
	LogFile             string   `short:"f" long:"file" value-name:"FILE" description:"Path to log file"`
	Pattern             []string `short:"p" long:"pattern" required:"true" value-name:"PAT" description:"Pattern to search for. If specified multiple, they will be treated together with the AND operator"`
	SuppressPattern     bool     `long:"suppress-pattern" description:"Suppress pattern display"`
	Exclude             []string `short:"E" long:"exclude" value-name:"PAT" description:"Pattern to exclude from matching. If specified multiple, they will be treated together with the AND operator"`
	WarnOver            int64    `short:"w" long:"warning-over" description:"Trigger a warning if matched lines is over a number"`
	CritOver            int64    `short:"c" long:"critical-over" description:"Trigger a critical if matched lines is over a number"`
	WarnLevel           float64  `long:"warning-level" value-name:"N" description:"Warning level if pattern has a group"`
	CritLevel           float64  `long:"critical-level" value-name:"N" description:"Critical level if pattern has a group"`
	ReturnContent       bool     `short:"r" long:"return" description:"Return matched line"`
	Directory           string   `long:"search-in-directory" value-name:"DIR" description:"Specify the directory of files to be detected"`
	FilePattern         string   `short:"F" long:"file-pattern" value-name:"FILE" description:"Check a pattern of files, instead of one file"`
	CaseInsensitive     bool     `short:"i" long:"icase" description:"Run a case insensitive match"`
	StateDir            string   `short:"s" long:"state-dir" value-name:"DIR" description:"Dir to keep state files under"`
	NoState             bool     `long:"no-state" description:"Don't use state file and read whole logs"`
	Encoding            string   `long:"encoding" description:"Encoding of log file"`
	Missing             string   `long:"missing" default:"UNKNOWN" value-name:"(CRITICAL|WARNING|OK|UNKNOWN)" description:"Exit status when log files missing"`
	CheckFirst          bool     `long:"check-first" description:"Check the log on the first run"`
	patternReg          []*regexp.Regexp
	excludeReg          []*regexp.Regexp
	fileListFromGlob    []string
	fileListFromPattern []string
	origArgs            []string
	decoder             *enc.Decoder

	testHookNewBufferedReader func(r io.Reader) *bufio.Reader
}

func (opts *logOpts) prepare() error {
	if opts.LogFile == "" && opts.FilePattern == "" {
		return fmt.Errorf("No log file specified")
	}

	if opts.Directory != "" && opts.FilePattern == "" {
		return fmt.Errorf("search-in-directory option must be used with file-pattern option")
	}

	var err error
	var reg *regexp.Regexp
	for _, ptn := range opts.Pattern {
		if reg, err = regCompileWithCase(ptn, opts.CaseInsensitive); err != nil {
			return fmt.Errorf("pattern is invalid")
		}
		opts.patternReg = append(opts.patternReg, reg)
	}

	if len(opts.patternReg) > 1 && (opts.WarnLevel > 0 || opts.CritLevel > 0) {
		return fmt.Errorf("When multiple patterns specified, --warning-level --critical-level can not be used")
	}

	for _, exclude := range opts.Exclude {
		if reg, err = regCompileWithCase(exclude, opts.CaseInsensitive); err != nil {
			return fmt.Errorf("exclude pattern is invalid")
		}
		opts.excludeReg = append(opts.excludeReg, reg)
	}

	if opts.LogFile != "" {
		opts.fileListFromGlob, err = zglob.Glob(opts.LogFile)
		// unless --missing specified, we should ignore file not found error
		if err != nil && err != os.ErrNotExist {
			return fmt.Errorf("invalid glob for --file")
		}
	}

	if opts.FilePattern != "" {
		opts.fileListFromPattern, err = parseFilePattern(opts.Directory, opts.FilePattern, opts.CaseInsensitive)
		if err != nil {
			return err
		}
	}
	if !validateMissing(opts.Missing) {
		return fmt.Errorf("missing option is invalid")
	}
	return nil
}

// Do the plugin
func Do() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	go func() {
		sig := <-sigCh
		log.Printf("check-log is exiting: caught a signal: %v", sig)
		cancel()
	}()
	signal.Notify(sigCh, defaultSignal)

	ckr := run(ctx, os.Args[1:])
	ckr.Name = "LOG"
	ckr.Exit()
}

func regCompileWithCase(ptn string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		ptn = "(?i)" + ptn
	}
	return regexp.Compile(ptn)
}

func validateMissing(missing string) bool {
	switch missing {
	case "CRITICAL", "WARNING", "OK", "UNKNOWN", "":
		return true
	default:
		return false
	}
}

func parseArgs(args []string) (*logOpts, error) {
	origArgs := make([]string, len(args))
	copy(origArgs, args)
	opts := &logOpts{}
	_, err := flags.ParseArgs(opts, args)
	opts.origArgs = origArgs
	if opts.StateDir == "" {
		workdir := pluginutil.PluginWorkDir()
		opts.StateDir = filepath.Join(workdir, "check-log")
	}
	return opts, err
}

func run(ctx context.Context, args []string) *checkers.Checker {
	opts, err := parseArgs(args)
	if err != nil {
		os.Exit(1)
	}

	err = opts.prepare()
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	warnNum := int64(0)
	critNum := int64(0)
	var missingFiles []string
	errorOverall := ""

	if opts.LogFile != "" && len(opts.fileListFromGlob) == 0 {
		missingFiles = append(missingFiles, opts.LogFile)
	}

	for _, f := range append(opts.fileListFromGlob, opts.fileListFromPattern...) {
		if ctx.Err() != nil {
			break
		}
		_, err := os.Stat(f)
		if err != nil {
			missingFiles = append(missingFiles, f)
			continue
		}
		w, c, errLines, err := opts.searchLog(ctx, f)
		if err != nil {
			return checkers.Unknown(err.Error())
		}
		warnNum += w
		critNum += c
		if opts.ReturnContent && errLines != "" {
			errorOverall += "[" + f + "]\n" + errLines
		}
	}

	var patterns []string
	for _, ptn := range opts.Pattern {
		patterns = append(patterns, fmt.Sprintf("/%s/", ptn))
	}
	var msg string
	if opts.SuppressPattern {
		msg = fmt.Sprintf("%d warnings, %d criticals.", warnNum, critNum)
	} else {
		msg = fmt.Sprintf("%d warnings, %d criticals for pattern %s.", warnNum, critNum, strings.Join(patterns, " and "))
	}
	if errorOverall != "" {
		msg += "\n" + errorOverall
	}
	checkSt := checkers.OK
	if len(missingFiles) > 0 {
		switch opts.Missing {
		case "OK":
		case "WARNING":
			checkSt = checkers.WARNING
		case "CRITICAL":
			checkSt = checkers.CRITICAL
		default:
			checkSt = checkers.UNKNOWN
		}
		msg += "\n" + fmt.Sprintf("The following %d files are missing.", len(missingFiles))
		for _, f := range missingFiles {
			msg += "\n" + f
		}
	}
	if warnNum > opts.WarnOver {
		checkSt = checkers.WARNING
	}
	if critNum > opts.CritOver {
		checkSt = checkers.CRITICAL
	}
	return checkers.NewChecker(checkSt, msg)
}

func (opts *logOpts) searchLog(ctx context.Context, logFile string) (int64, int64, string, error) {
	if ctx.Err() != nil {
		return 0, 0, "", nil
	}
	stateFile := getStateFile(opts.StateDir, logFile, opts.origArgs)
	skipBytes, inode, isFirstCheck := int64(0), uint(0), false
	if !opts.NoState {
		s, err := getBytesToSkip(stateFile)
		if err != nil {
			if err != errValidStateFileNotFound {
				return 0, 0, "", err
			}
			isFirstCheck = true
		}
		skipBytes = s

		i, err := getInode(stateFile)
		if err != nil {
			return 0, 0, "", err
		}
		inode = i
	}

	f, err := os.Open(logFile)
	if err != nil {
		return 0, 0, "", err
	}
	defer f.Close()

	var oldf *os.File
	if !opts.NoState {
		oldf, err = openOldFile(logFile, &state{SkipBytes: skipBytes, Inode: inode})
		if err != nil {
			return 0, 0, "", err
		}
		defer oldf.Close()
	}

	stat, err := f.Stat()
	if err != nil {
		return 0, 0, "", err
	}

	// Skip whole file on first check, unless CheckFirst specified
	if !opts.NoState && isFirstCheck && !opts.CheckFirst {
		skipBytes = stat.Size()
	}

	fmt.Println("=================================")
	fmt.Println("inode-> ", inode)
	fmt.Println("detectInode inode-> ", detectInode(stat))
	rotated := false
	if stat.Size() < skipBytes || (inode != 0 && inode != detectInode(stat)) {
		rotated = true
	} else if skipBytes > 0 {
		f.Seek(skipBytes, 0)
	}

	var r io.Reader = f
	var oldr io.Reader = oldf
	if opts.Encoding != "" {
		e := encoding.GetEncoding(opts.Encoding)
		if e == nil {
			return 0, 0, "", fmt.Errorf("unknown encoding:" + opts.Encoding)
		}
		opts.decoder = e.NewDecoder()
	}

	warnNum, critNum, readBytes, errLines, err := opts.searchReader(ctx, r)
	if err != nil {
		return warnNum, critNum, errLines, err
	}

	if oldf != nil {
		// search old file
		var (
			oldWarnNum, oldCritNum int64
			oldErrLines            string
		)
		// ignore readBytes under the premise that the old file will never be updated.
		oldWarnNum, oldCritNum, _, oldErrLines, err := opts.searchReader(ctx, oldr)
		if err != nil {
			return oldWarnNum, critNum, errLines, err
		}
		warnNum += oldWarnNum
		critNum += oldCritNum
		errLines += oldErrLines
	}

	if rotated {
		skipBytes = readBytes
	} else {
		skipBytes += readBytes
	}

	if !opts.NoState {
		err = saveState(stateFile, &state{SkipBytes: skipBytes, Inode: detectInode(stat)})
		if err != nil {
			log.Printf("writeByteToSkip failed: %s\n", err.Error())
		}
	}
	return warnNum, critNum, errLines, nil
}

func newBufferedReader(r io.Reader) *bufio.Reader {
	return bufio.NewReader(r)
}

func (opts *logOpts) searchReader(ctx context.Context, rdr io.Reader) (warnNum, critNum, readBytes int64, errLines string, err error) {
	newReader := opts.testHookNewBufferedReader
	if newReader == nil {
		newReader = newBufferedReader
	}

	var errLinesBuilder strings.Builder
	r := newReader(rdr)
	for ctx.Err() == nil {
		lineBytes, rErr := r.ReadBytes('\n')
		if rErr != nil {
			if rErr != io.EOF {
				err = rErr
			}
			break
		}
		readBytes += int64(len(lineBytes))

		if opts.decoder != nil {
			lineBytes, err = opts.decoder.Bytes(lineBytes)
			if err != nil {
				break
			}
		}
		line := strings.Trim(string(lineBytes), "\r\n")
		if matched, matches := opts.match(line); matched {
			if len(matches) > 1 && (opts.WarnLevel > 0 || opts.CritLevel > 0) {
				level, err := strconv.ParseFloat(matches[1], 64)
				if err != nil {
					warnNum++
					critNum++
					errLinesBuilder.WriteString(line)
					errLinesBuilder.WriteString("\n")
				} else {
					levelOver := false
					if level > opts.WarnLevel {
						levelOver = true
						warnNum++
					}
					if level > opts.CritLevel {
						levelOver = true
						critNum++
					}
					if levelOver {
						errLinesBuilder.WriteString(line)
						errLinesBuilder.WriteString("\n")
					}
				}
			} else {
				warnNum++
				critNum++
				errLinesBuilder.WriteString(line)
				errLinesBuilder.WriteString("\n")
			}
		}
	}

	errLines = errLinesBuilder.String()
	return
}

func (opts *logOpts) match(line string) (bool, []string) {
	var matches []string
	for _, pReg := range opts.patternReg {
		matches = pReg.FindStringSubmatch(line)
		if len(matches) == 0 {
			return false, nil
		}
	}
	if len(opts.excludeReg) > 0 {
		exclude := true
		for _, eReg := range opts.excludeReg {
			if !eReg.MatchString(line) {
				exclude = false
				break
			}
		}
		if exclude {
			return false, nil
		}
	}
	return true, matches
}

func parseFilePattern(directory, filePattern string, caseInsensitive bool) ([]string, error) {
	var dirStr string
	var filePat string
	if directory != "" {
		dirStr = directory
		filePat = filePattern
	} else {
		// for backward compatibility.
		// Directory delimiters and regular expressions can conflict in windows environment.
		dirStr = filepath.Dir(filePattern)
		filePat = filepath.Base(filePattern)
	}
	reg, err := regCompileWithCase(filePat, caseInsensitive)
	if err != nil {
		return nil, fmt.Errorf("file-pattern is invalid")
	}

	fileInfos, err := os.ReadDir(dirStr)
	if err != nil {
		return nil, fmt.Errorf("cannot read the directory:" + err.Error())
	}

	var fileList []string
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		fname := fileInfo.Name()
		if reg.MatchString(fname) {
			fileList = append(fileList, dirStr+string(filepath.Separator)+fileInfo.Name())
		}
	}
	return fileList, nil
}

type state struct {
	SkipBytes int64 `json:"skip_bytes"`
	Inode     uint  `json:"inode"`
}

func loadState(fname string) (*state, error) {
	state := &state{}
	b, err := os.ReadFile(fname)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return state, err
	}
	err = json.Unmarshal(b, state)
	if err != nil {
		// this json unmarshal error will be ignored by callers
		log.Printf("failed to loadState (will be ignored): %s", err)
		return nil, errStateFileCorrupted
	}
	return state, nil
}

var stateRe = regexp.MustCompile(`^([a-zA-Z]):[/\\]`)

func getStateFile(stateDir, f string, args []string) string {
	return filepath.Join(
		stateDir,
		fmt.Sprintf(
			"%s-%x.json",
			stateRe.ReplaceAllString(f, `$1`+string(filepath.Separator)),
			md5.Sum([]byte(strings.Join(args, " "))),
		),
	)
}

var errValidStateFileNotFound = fmt.Errorf("state file not found, or corrupted")
var errStateFileCorrupted = fmt.Errorf("state file is corrupted")

func getBytesToSkip(f string) (int64, error) {
	state, err := loadState(f)
	// Do not fallback to old status file when JSON file is corrupted
	if err == errStateFileCorrupted {
		return 0, errValidStateFileNotFound
	}
	if err != nil {
		return 0, err
	}
	if state != nil {
		// json file exists
		return state.SkipBytes, nil
	}
	// Fallback to read old style status file
	// for backward compatibility.
	// Once saved as new style file, the following will be unreachable.
	oldf := strings.TrimSuffix(f, ".json")
	return getBytesToSkipOld(oldf)
}

func getBytesToSkipOld(f string) (int64, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, errValidStateFileNotFound
		}
		return 0, err
	}

	i, err := strconv.ParseInt(strings.Trim(string(b), " \r\n"), 10, 64)
	if err != nil {
		log.Printf("failed to getBytesToSkip (ignoring): %s", err)
	}
	return i, nil
}

func getInode(f string) (uint, error) {
	state, err := loadState(f)
	// ignore corrupted json
	if err == errStateFileCorrupted {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	if state != nil {
		// json file exists
		return state.Inode, nil
	}
	return 0, nil
}

func saveState(f string, state *state) error {
	b, _ := json.Marshal(state)
	if err := os.MkdirAll(filepath.Dir(f), 0755); err != nil {
		return err
	}
	return atomic.WriteFile(f, bytes.NewReader(b))
}

var errFileNotFoundByInode = fmt.Errorf("old file not found")

func findFileByInode(inode uint, dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		fi, err := entry.Info()
		if err != nil {
			return "", err
		}
		if detectInode(fi) == inode {
			return filepath.Join(dir, fi.Name()), nil
		}
	}
	return "", errFileNotFoundByInode
}

func openOldFile(f string, state *state) (*os.File, error) {
	fi, err := os.Stat(f)
	if err != nil {
		return nil, err
	}
	inode := detectInode(fi)
	if state.Inode > 0 && state.Inode != inode {
		if oldFile, err := findFileByInode(state.Inode, filepath.Dir(f)); err == nil {
			oldf, err := os.Open(oldFile)
			if err != nil {
				return nil, err
			}
			oldfi, _ := oldf.Stat()
			if oldfi.Size() > state.SkipBytes {
				oldf.Seek(state.SkipBytes, io.SeekStart)
				return oldf, nil
			}
		} else if err != errFileNotFoundByInode {
			return nil, err
		}
		// just ignore the process of searching old file if errFileNotFoundByInode
	}
	return nil, nil
}
