package checkuptime

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/mackerelio/golib/uptime"
)

var opts struct {
	WarnUnder    *float64 `long:"warn-under" value-name:"N" description:"(DEPRECATED) Trigger a warning if under the seconds"`
	WarningUnder *float64 `short:"w" long:"warning-under" value-name:"N" description:"Trigger a warning if under the seconds"`
	CritUnder    *float64 `short:"c" long:"critical-under" value-name:"N" description:"Trigger a critial if under the seconds"`
	WarnOver     *float64 `long:"warn-over" value-name:"N" description:"(DEPRECATED) Trigger a warning if over the seconds"`
	WarningOver  *float64 `short:"W" long:"warning-over" value-name:"N" description:"Trigger a warning if over the seconds"`
	CritOver     *float64 `short:"C" long:"critical-over" value-name:"N" description:"Trigger a critical if over the seconds"`
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "Uptime"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}
	ut, err := uptime.Get()
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Faild to fetch uptime metrics: %s", err))
	}

	// for backward compatibility
	if opts.WarnUnder != nil && opts.WarningUnder == nil {
		opts.WarningUnder = opts.WarnUnder
	}
	if opts.WarnOver != nil && opts.WarningOver == nil {
		opts.WarningOver = opts.WarnOver
	}

	checkSt := checkers.OK
	if opts.WarningUnder != nil && *opts.WarningUnder > ut {
		checkSt = checkers.WARNING
	}
	if opts.WarningOver != nil && *opts.WarningOver < ut {
		checkSt = checkers.WARNING
	}
	if opts.CritUnder != nil && *opts.CritUnder > ut {
		checkSt = checkers.CRITICAL
	}
	if opts.CritOver != nil && *opts.CritOver < ut {
		checkSt = checkers.CRITICAL
	}
	dur := time.Duration(ut * float64(time.Second))
	hours := int64(dur.Hours())
	days := hours / 24
	hours = hours % 24
	mins := int64(dur.Minutes()) % 60
	msg := fmt.Sprintf("%d day(s) %d hour(s) %d minute(s) (%d second(s))\n", days, hours, mins, int64(dur.Seconds()))

	return checkers.NewChecker(checkSt, msg)
}
