package checkntpoffset

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var opts struct {
	Crit         float64 `short:"c" long:"critical" default:"100" description:"Critical threshold of ntp offset(ms)"`
	Warn         float64 `short:"w" long:"warning" default:"50" description:"Warning threshold of ntp offset(ms)"`
	NTPServers   string  `short:"s" long:"ntp-servers" default:"" description:"Use specified NTP Servers(plural servers can be set separated by ,). When set plural servers, use first response. If not set, use local command just like ntpd/chronyd."`
	NTPTimeout   int     `short:"t" long:"ntp-timeout" default:"15" description:"Timeout of NTP Server Querying(in seconds)."`
	CheckStratum bool    `short:"S" long:"check-stratum" description:"Check stratum and fail if the machine is not synchronized."`
}

var ntpTimeout int

func init() {
	ntpTimeout = 15
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
	ntpTimeout = opts.NTPTimeout

	offset, err := getNTPOffset(opts.NTPServers, opts.CheckStratum)
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
		msg = fmt.Sprintf("ntp offset is %f(actual) < %f(warning threshold), %f(critial threshold)", math.Abs(offset), opts.Warn, opts.Crit)
		chkSt = checkers.OK
	}

	return checkers.NewChecker(chkSt, msg)
}

func withCmd(cmd *exec.Cmd, fn func(io.Reader) error) error {
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := fn(out); err != nil {
		return err
	}
	return cmd.Wait()
}

const (
	ntpNTPD    = "ntpd"
	ntpChronyd = "chronyd"

	cmdNTPq    = "ntpq"
	cmdChronyc = "chronyc"
)

func hasCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func detectNTPDname() (ntpdName string, err error) {
	err = withCmd(exec.Command("ps", "-eo", "comm"), func(out io.Reader) error {
		scr := bufio.NewScanner(out)
		for scr.Scan() {
			switch filepath.Base(scr.Text()) {
			case ntpChronyd:
				if hasCommand(cmdChronyc) {
					ntpdName = ntpChronyd
					return nil
				}
			case ntpNTPD:
				if hasCommand(cmdNTPq) {
					ntpdName = ntpNTPD
					return nil
				}
			}
		}
		return fmt.Errorf("no ntp daemons detected")
	})
	return ntpdName, err
}

func getNTPOffset(ntpServers string, checkStratum bool) (float64, error) {
	if ntpServers != "" {
		return getNTPOffsetFromNTPServers(ntpServers)
	}

	ntpdName, err := detectNTPDname()
	if err != nil {
		return 0.0, err
	}
	switch ntpdName {
	case ntpNTPD:
		return getNTPOffsetFromNTPD(checkStratum)
	case ntpChronyd:
		return getNTPOffsetFromChrony(checkStratum)
	}
	return 0.0, fmt.Errorf("unsupported ntp daemon %q", ntpdName)
}

// getNTPOffsetFromNTPServers ask time to ntp servers and return NTP Offset.
// Use first response, ignore others
//
// FIXME need fluent cancel mechanism
func getNTPOffsetFromNTPServers(ntpServers string) (offset float64, err error) {
	resultChan := make(chan float64)
	for _, ntpServer := range strings.Split(ntpServers, ",") {
		go func(ntpServer string) error {
			ntpServer = strings.Trim(ntpServer, " ")
			options := ntp.QueryOptions{Timeout: time.Duration(ntpTimeout) * time.Second}
			response, err := ntp.QueryWithOptions(ntpServer, options)
			if err != nil {
				return err
			}
			resultChan <- float64(response.ClockOffset / time.Millisecond)
			return nil
		}(ntpServer)
	}

	select {
	case <-time.After(time.Duration(ntpTimeout) * time.Second):
		// return error only when all NTPServers are failed
		return 0.0, fmt.Errorf("NTP offset cannot get from %q", ntpServers)
	case offset = <-resultChan:
		return offset, nil
	}
}

func getNTPOffsetFromNTPD(checkStratum bool) (offset float64, err error) {
	err = withCmd(exec.Command(cmdNTPq, "-c", "rv 0 stratum,offset"), func(out io.Reader) error {
		offset, err = parseNTPOffsetFromNTPD(out, checkStratum)
		return err
	})
	return offset, err
}

func parseNTPOffsetFromNTPD(out io.Reader, checkStratum bool) (float64, error) {
	scr := bufio.NewScanner(out)
	const stratumPrefix = "stratum="
	const offsetPrefix = "offset="
	var offset *float64
	for scr.Scan() {
		line := scr.Text()
		for _, column := range strings.Split(line, ",") {
			column = strings.TrimPrefix(column, " ")
			if checkStratum && strings.HasPrefix(column, stratumPrefix) {
				stratum, err := strconv.ParseInt(strings.TrimPrefix(column, stratumPrefix), 10, 64)
				if err != nil {
					return 0.0, err
				}
				// stratum == 16 means that the machine is unsynchronized.
				// ref. https://support.ntp.org/bin/view/Support/TroubleshootingNTP
				if stratum == 16 {
					return 0.0, fmt.Errorf("not synchronized to stratum")
				}
			}
			if offset == nil && strings.HasPrefix(column, offsetPrefix) {
				offsetMillis, err := strconv.ParseFloat(strings.TrimPrefix(column, offsetPrefix), 64)
				if err != nil {
					return 0.0, err
				}
				offset = &offsetMillis
			}
		}
	}
	if offset == nil {
		return 0.0, fmt.Errorf("couldn't get ntp offset. ntpd process may be down")
	}
	return *offset, nil
}

func getNTPOffsetFromChrony(checkStratum bool) (offset float64, err error) {
	err = withCmd(exec.Command(cmdChronyc, "tracking"), func(out io.Reader) error {
		offset, err = parseNTPOffsetFromChrony(out, checkStratum)
		return err
	})
	return offset, err
}

func parseNTPOffsetFromChrony(out io.Reader, checkStratum bool) (float64, error) {
	scr := bufio.NewScanner(out)
	const stratumPrefix = "Stratum"
	const offsetPrefix = "Last offset"
	var offset *float64
	for scr.Scan() {
		line := scr.Text()
		if checkStratum && strings.HasPrefix(line, stratumPrefix) {
			flds := strings.Fields(line)
			if len(flds) != 3 {
				return 0.0, fmt.Errorf("failed to get ntp stratum")
			}
			stratum, err := strconv.ParseInt(flds[2], 10, 64)
			if err != nil {
				return 0.0, err
			}
			// stratum == 0 means that the machine is unsynchronized.
			// Actually this can be the best case, but that would be rare...
			if stratum == 0 {
				return 0.0, fmt.Errorf("not synchronized to stratum")
			}
		}
		if offset == nil && strings.HasPrefix(line, offsetPrefix) {
			flds := strings.Fields(line)
			if len(flds) != 5 {
				return 0.0, fmt.Errorf("failed to get ntp offset")
			}
			offsetSeconds, err := strconv.ParseFloat(flds[3], 64)
			if err != nil {
				return 0.0, err
			}
			offsetMillis := offsetSeconds * 1000
			offset = &offsetMillis
		}
	}
	if offset == nil {
		return 0.0, fmt.Errorf("failed to get ntp offset")
	}
	return *offset, nil
}
