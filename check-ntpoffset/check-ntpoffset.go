package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var opts struct {
	Crit float64 `short:"c" long:"critical" default:"100" description:"critical if the ntpoffset is over"`
	Warn float64 `short:"w" long:"warning" default:"50" description:"warning if the ntpoffset is over"`
}

func main() {
	ckr := run(os.Args[1:])
	ckr.Name = "NTP"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return checkers.Critical(err.Error())
	}

	offset, err := getNtpOffset()
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	var chkSt checkers.Status
	var msg string
	if opts.Crit < math.Abs(offset) {
		msg = fmt.Sprintf("ntp offset is over %f(actual) > %f(threshold)", math.Abs(offset), opts.Crit)
		chkSt = checkers.CRITICAL
	} else if opts.Warn < math.Abs(offset) {
		msg = fmt.Sprintf("ntp offset is over %f(actual) > %f(threshold)", math.Abs(offset), opts.Warn)
		chkSt = checkers.WARNING
	} else {
		msg = fmt.Sprintf("ntp offset is %f(actual) < %f(warning threshold), %f(critial threshold)", offset, opts.Warn, opts.Crit)
		chkSt = checkers.OK
	}

	return checkers.NewChecker(chkSt, msg)
}

func getNtpOffset() (float64, error) {
	var offset float64
	var err error

	output, err := exec.Command("ntpq", "-c", "rv 0 offset").Output()
	if err != nil {
		return offset, err
	}

	o := strings.Split(string(output), "=")
	if len(o) != 2 {
		return offset, errors.New("couldn't get ntp offset. ntpd process may be down")
	}

	offset, err = strconv.ParseFloat(strings.Trim(o[1], "\n"), 64)
	if err != nil {
		return offset, err
	}

	return offset, nil
}
