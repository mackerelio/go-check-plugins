package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type logOpts struct {
	StateDir        string  `short:"s" long:"state-dir" default:"/var/mackerel-cache/check-log" value-name:"DIR" description:"Dir to keep state files under"`
	LogFile         string  `short:"f" long:"log-file" value-name:"FILE" description:"Path to log file"`
	Pattern         string  `short:"p" long:"pattern" required:"true" value-name:"PAT" description:"Pattern to search for"`
	Exclude         string  `short:"E" long:"exclude" value-name:"PAT" description:"Pattern to exclude from matching"`
	WarnOver        int64   `short:"w" long:"warning-over" description:"Trigger a warning if matched lines is over a number"`
	CritOver        int64   `short:"c" long:"critical-over" description:"Trigger a critical if matched lines is over a number"`
	WarnLevel       float64 `long:"warning-level" value-name:"N" description:"Warning level if pattern has a group"`
	CritLevel       float64 `long:"critical-level" value-name:"N" description:"Critical level if pattern has a group"`
	CaseInsensitive bool    `short:"i" long:"icase" description:"Run a case insensitive match"`
	FilePattern     string  `short:"F" long:"filepattern" value-name:"FILE" description:"Check a pattern of files, instead of one file"`
	ReturnContent   bool    `short:"r" long:"return" description:"Return matched line"`
	patternReg      *regexp.Regexp
	excludeReg      *regexp.Regexp
	fileList        []string
}

func (opts *logOpts) prepare() error {
	if opts.LogFile == "" && opts.FilePattern == "" {
		return fmt.Errorf("No log file specified")
	}

	var err error
	if opts.patternReg, err = regCompileWithCase(opts.Pattern, opts.CaseInsensitive); err != nil {
		return fmt.Errorf("pattern is invalid")
	}

	if opts.Exclude != "" {
		opts.excludeReg, err = regCompileWithCase(opts.Exclude, opts.CaseInsensitive)
		if err != nil {
			return fmt.Errorf("exclude pattern is invalid")
		}
	}

	if opts.LogFile != "" {
		opts.fileList = append(opts.fileList, opts.LogFile)
	}

	if opts.FilePattern != "" {
		dirStr := filepath.Dir(opts.FilePattern)
		filePat := filepath.Base(opts.FilePattern)
		reg, err := regCompileWithCase(filePat, opts.CaseInsensitive)
		if err != nil {
			return fmt.Errorf("file-pattern is invalid")
		}

		fileInfos, err := ioutil.ReadDir(dirStr)
		if err != nil {
			return fmt.Errorf("cannot read the Directory:" + err.Error())
		}

		for _, fileInfo := range fileInfos {
			if fileInfo.IsDir() {
				continue
			}
			fname := fileInfo.Name()
			if opts.CaseInsensitive {
				fname = strings.ToLower(fname)
			}
			if reg.MatchString(fname) {
				opts.fileList = append(opts.fileList, dirStr+string(filepath.Separator)+fileInfo.Name())
			}
		}
	}
	return nil
}

func main() {
	ckr := run(os.Args[1:])
	ckr.Name = "LOG"
	ckr.Exit()
}

func regCompileWithCase(ptn string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		ptn = strings.ToLower(ptn)
	}
	return regexp.Compile(ptn)
}

func run(args []string) *checkers.Checker {
	opts := &logOpts{}
	_, err := flags.ParseArgs(opts, args)
	if err != nil {
		os.Exit(1)
	}

	err = opts.prepare()
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	warnNum := int64(0)
	critNum := int64(0)
	errorOverall := ""

	for _, f := range opts.fileList {
		w, c, errLines, err := opts.searchLog(f)
		if err != nil {
			return checkers.Unknown(err.Error())
		}
		warnNum += w
		critNum += c
		if opts.ReturnContent {
			errorOverall += errLines
		}
	}

	checkSt := checkers.OK
	if warnNum > opts.WarnOver {
		checkSt = checkers.WARNING
	}
	if critNum > opts.CritOver {
		checkSt = checkers.CRITICAL
	}
	msg := fmt.Sprintf("%d warnings, %d criticals for pattern %s. %s", warnNum, critNum, opts.Pattern, errorOverall)
	return checkers.NewChecker(checkSt, msg)
}

func (opts *logOpts) searchLog(logFile string) (int64, int64, string, error) {
	stateFile := getStateFile(opts.StateDir, logFile)
	skipBytes, err := getBytesToSkip(stateFile)
	if err != nil {
		return 0, 0, "", err
	}
	f, err := os.Open(logFile)
	if err != nil {
		return 0, 0, "", err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return 0, 0, "", err
	}

	if skipBytes > 0 && stat.Size() >= skipBytes {
		f.Seek(skipBytes, 0)
	}

	warnNum, critNum, readBytes, errLines, err := opts.searchReader(f)

	err = writeBytesToSkip(stateFile, readBytes+skipBytes)
	if err != nil {
		log.Printf("writeByteToSkip failed: %s\n", err.Error())
	}
	return warnNum, critNum, errLines, nil
}

func (opts *logOpts) searchReader(rdr io.Reader) (warnNum, critNum, readBytes int64, errLines string, err error) {
	r := bufio.NewReader(rdr)
	for {
		lineBytes, rErr := r.ReadBytes('\n')
		readBytes += int64(len(lineBytes))
		if rErr != nil {
			if rErr != io.EOF {
				err = rErr
			}
			break
		}
		line := strings.Trim(string(lineBytes), "\r\n")
		checkLine := line
		if opts.CaseInsensitive {
			checkLine = strings.ToLower(checkLine)
		}
		if matched, matches := opts.match(checkLine); matched {
			if len(matches) > 1 && (opts.WarnLevel > 0 || opts.CritLevel > 0) {
				level, err := strconv.ParseFloat(matches[1], 64)
				if err != nil {
					warnNum++
					critNum++
					errLines += "\n" + line
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
						errLines += "\n" + line
					}
				}
			} else {
				warnNum++
				critNum++
				errLines += "\n" + line
			}
		}
	}
	return
}

func (opts *logOpts) match(line string) (bool, []string) {
	pReg := opts.patternReg
	eReg := opts.excludeReg

	matches := pReg.FindStringSubmatch(line)
	matched := len(matches) > 0 && (eReg == nil || !eReg.MatchString(line))
	return matched, matches
}

var stateRe = regexp.MustCompile(`^([A-Z]):[/\\]`)

func getStateFile(stateDir, f string) string {
	return filepath.Join(stateDir, stateRe.ReplaceAllString(f, `$1`+string(filepath.Separator)))
}

func getBytesToSkip(f string) (int64, error) {
	_, err := os.Stat(f)
	if err != nil {
		return 0, nil
	}
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(strings.Trim(string(b), " \r\n"))
	if err != nil {
		return 0, err
	}
	return int64(i), nil
}

func writeBytesToSkip(f string, num int64) error {
	err := os.MkdirAll(filepath.Dir(f), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, []byte(fmt.Sprintf("%d", num)), 0755)
}
