package checkfilesize

import (
	"fmt"
	"os"

	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "FileSize"
	ckr.Exit()
}

var opts struct {
	Base  string `short:"b" long:"base" required:"true" description:"base directory"`
	Warn  string `short:"w" long:"warning" default:"1K" description:"warning if the size is over"`
	Crit  string `short:"c" long:"critical" default:"1K" description:"critical if the size is over"`
	Depth int    `short:"d" long:"depth" default:"1" description:"max depth of the directory from base directory"`
}

var sizeReg = regexp.MustCompile(`^(\d+\.?\d*)(k|K|m|M|g|G|t|T)?$`)

func sizeValue(input string) (float64, error) {
	var size float64
	var err error

	r := sizeReg.FindAllStringSubmatch(input, -1)
	if len(r) == 0 {
		msg := fmt.Errorf("%s is invalid", input)
		return -1, msg
	}
	num := r[0][1]
	unit := r[0][2]

	size, err = strconv.ParseFloat(num, 10)
	if err != nil {
		return -1, err
	}
	switch strings.ToLower(unit) {
	case "k":
		size = size * 1024
	case "m":
		size = size * 1024 * 1024
	case "g":
		size = size * 1024 * 1024 * 1024
	case "t":
		size = size * 1024 * 1024 * 1024 * 1024
	}

	return size, err
}

func listFiles(base string, depth int) ([]string, error) {

	var files []string

	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(base, path)
		if len(strings.Split(rel, "/")) > depth {
			return nil
		}
		files = append(files, path)

		return err
	})

	return files, err
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	ws, err := sizeValue(opts.Warn)
	if err != nil {
		return checkers.NewChecker(checkers.UNKNOWN, err.Error())
	}

	cs, err := sizeValue(opts.Crit)
	if err != nil {
		return checkers.NewChecker(checkers.UNKNOWN, err.Error())
	}

	var stat os.FileInfo
	var size int64
	files, err := listFiles(opts.Base, opts.Depth)
	if err != nil {
		// file not found
		return checkers.NewChecker(checkers.OK, err.Error())
	}

	for _, v := range files {
		stat, err = os.Stat(v)
		if err != nil {
			continue
		}
		size = size + stat.Size()
	}

	var chkSt checkers.Status
	var msg string
	if cs < float64(size) {
		msg = fmt.Sprintf("size %d Byte > %s Byte in %s", size, opts.Crit, opts.Base)
		chkSt = checkers.CRITICAL
	} else if ws < float64(size) {
		msg = fmt.Sprintf("size %d Byte > %s Byte in %s", size, opts.Warn, opts.Base)
		chkSt = checkers.WARNING
	} else {
		msg = fmt.Sprintf("size %d Byte < warning %s Byte, critical %s Byte in %s", size, opts.Warn, opts.Crit, opts.Base)
		chkSt = checkers.OK
	}

	return checkers.NewChecker(chkSt, msg)
}
