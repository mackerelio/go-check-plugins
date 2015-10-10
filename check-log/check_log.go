package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var opts struct {
	StateDir        string `short:"s" long:"state-dir" default:"/var/mackerel-cache/check-log" value-name:"DIR" description:"Dir to keep state files under"`
	LogFile         string `short:"f" long:"log-file" value-name:"FILE" description:"Path to log file"`
	Pattern         string `short:"q" long:"pattern" required:"true" value-name:"PAT" description:"Pattern to search for"`
	Exclude         string `short:"E" long:"exclude" default:"(?!)" value-name:"PAT" description:"Pattern to exclude from matching"`
	Warn            int64  `short:"w" long:"warn" value-name:"N" description:"Warning level if pattern has a group"`
	Crit            int64  `short:"c" long:"crit" value-name:"N" description:"Critical level if pattern has a group"`
	OnlyWarn        bool   `short:"o" long:"warn-only" description:"Warn instead of critical on match"`
	CaseInsensitive bool   `short:"i" long:"icase" description:"Run a case insensitive match"`
	FilePattern     string `short:"F" long:"filepattern" value-name:"FILE" description:"Check a pattern of files, instead of one file"`
	ReturnContent   bool   `short:"r" long:"return" description:"Return matched line"`
}

func main() {
	ckr := run(os.Args[1:])
	ckr.Name = "LOG"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	if opts.LogFile == "" && opts.FilePattern == "" {
		return checkers.Unknown("No log file specified")
	}

	excludeReg, err := regexp.Compile(opts.Exclude)
	if err != nil {
		return checkers.Unknown("exclude pattern is invalid")
	}
	excludeReg = excludeReg

	fileList := []string{}
	if opts.LogFile != "" {
		fileList = append(fileList, opts.LogFile)
	}

	if opts.FilePattern != "" {
		dirStr := filepath.Dir(opts.FilePattern)
		filePat := filepath.Base(opts.FilePattern)
		if opts.CaseInsensitive {
			filePat = strings.ToLower(filePat)
		}
		reg, err := regexp.Compile(filePat)
		if err != nil {
			return checkers.Unknown("file-pattern is invalid")
		}

		fileInfos, err := ioutil.ReadDir(dirStr)
		if err != nil {
			return checkers.Unknown("cannot read the Directory:" + err.Error())
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
				fileList = append(fileList, dirStr+string(filepath.Separator)+fileInfo.Name())
			}
		}
	}

	warnNum := 0
	critNum := 0
	errorOverall := ""

	// for _, _ = range fileList {
	//}

	checkSt := checkers.OK
	if warnNum > 0 {
		checkSt = checkers.WARNING
	}
	if critNum > 0 {
		checkSt = checkers.CRITICAL
	}
	msg := fmt.Sprintf("%d warnings, %d criticals for pattern %s. %s", warnNum, critNum, opts.Pattern, errorOverall)
	return checkers.NewChecker(checkSt, msg)
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
