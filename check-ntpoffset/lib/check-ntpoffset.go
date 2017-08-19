package checkntpoffset

import (
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
	Crit float64 `short:"c" long:"critical" default:"100" description:"Critical threshold of ntp offset(ms)"`
	Warn float64 `short:"w" long:"warning" default:"50" description:"Warning threshold of ntp offset(ms)"`
}

// Do the plugin
func Do() {
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

func getNtpOffset() (offset float64, err error) {
	psout, err := exec.Command("ps", "-eo", "comm").Output()
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(psout), "\n") {
		if strings.HasSuffix(line, "chronyd") {
			return getChronyNtpOffset()
		}
	}

	output, err := exec.Command("ntpq", "-c", "rv 0 offset").Output()
	if err != nil {
		return
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
		return offset, fmt.Errorf("couldn't get ntp offset. ntpd process may be down")
	}

	o := strings.Split(string(line), "=")
	if len(o) != 2 {
		return offset, fmt.Errorf("couldn't get ntp offset. ntpd process may be down")
	}
	return strconv.ParseFloat(strings.Trim(o[1], "\n"), 64)
}

func getChronyNtpOffset() (offset float64, err error) {
	output, err := exec.Command("chronyc", "tracking").Output()
	// Reference ID    : 160.16.75.242 (sv01.azsx.net)
	// Stratum         : 3
	// Ref time (UTC)  : Thu May  4 11:51:30 2017
	// System time     : 0.000033190 seconds slow of NTP time
	// Last offset     : +0.000003614 seconds
	// RMS offset      : 0.000017540 seconds
	// Frequency       : 10.880 ppm fast
	// Residual freq   : -0.000 ppm
	// Skew            : 0.003 ppm
	// Root delay      : 0.003541 seconds
	// Root dispersion : 0.000849 seconds
	// Update interval : 1030.4 seconds
	// Leap status     : Normal
	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "Last offset") {
			flds := strings.Fields(line)
			if len(flds) != 5 {
				return 0.0, fmt.Errorf("failed to get ntp offset")
			}
			offset, err := strconv.ParseFloat(flds[3], 64)
			if err != nil {
				return 0.0, err
			}
			return offset * 1000, nil
		}
	}
	return 0.0, fmt.Errorf("failed to get ntp offset")
}
