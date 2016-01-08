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
		os.Exit(1)
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

	var line string
	lines := strings.Split(string(output), "\n")
	switch len(lines) {
	case 2:
		line = lines[0]
	case 3:
		/* example on ntp 4.2.2p1-18.el5.centos
		   assID=0 status=06f4 leap_none, sync_ntp, 15 events, event_peer/strat_chg,
		   offset=0.180
		*/

		if strings.Index(lines[0], `assID=0`) == 0 {
			line = lines[1]
			break
		}
		fallthrough
	default:
		return offset, errors.New("couldn't get ntp offset. ntpd process may be down")
	}

	o := strings.Split(string(line), "=")
	if len(o) != 2 {
		return offset, errors.New("couldn't get ntp offset. ntpd process may be down")
	}

	offset, err = strconv.ParseFloat(strings.Trim(o[1], "\n"), 64)
	if err != nil {
		return offset, err
	}

	return offset, nil
}
