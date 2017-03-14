package checkfileage

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "FileAge"
	ckr.Exit()
}

type monitor struct {
	warningAge   int64
	warningSize  int64
	criticalAge  int64
	criticalSize int64
}

func (m monitor) hasWarningAge() bool {
	return m.warningAge != 0
}

func (m monitor) hasWarningSize() bool {
	return m.warningSize != 0
}

func (m monitor) CheckWarning(age, size int64) bool {
	return (m.hasWarningAge() && m.warningAge < age) ||
		(m.hasWarningSize() && m.warningSize > size)
}

func (m monitor) hasCriticalAge() bool {
	return m.criticalAge != 0
}

func (m monitor) hasCriticalSize() bool {
	return m.criticalSize != 0
}

func (m monitor) CheckCritical(age, size int64) bool {
	return (m.hasCriticalAge() && m.criticalAge < age) ||
		(m.hasCriticalSize() && m.criticalSize > size)
}

func newMonitor(warningAge, warningSize, criticalAge, criticalSize int64) *monitor {
	return &monitor{
		warningAge:   warningAge,
		warningSize:  warningSize,
		criticalAge:  criticalAge,
		criticalSize: criticalSize,
	}
}

var opts struct {
	File          string `short:"f" long:"file" required:"true" description:"monitor file name"`
	WarningAge    int64  `short:"w" long:"warning-age" default:"240" description:"warning if more old than"`
	WarningSize   int64  `short:"W" long:"warning-size" description:"warning if file size less than"`
	CriticalAge   int64  `short:"c" long:"critical-age" default:"600" description:"critical if more old than"`
	CriticalSize  int64  `short:"C" long:"critical-size" default:"0" description:"critical if file size less than"`
	IgnoreMissing bool   `short:"i" long:"ignore-missing" description:"skip alert if file doesn't exist"`
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	stat, err := os.Stat(opts.File)
	if err != nil {
		if opts.IgnoreMissing {
			return checkers.Ok("No such file, but ignore missing is set.")
		}
		return checkers.Unknown(err.Error())
	}

	monitor := newMonitor(opts.WarningAge, opts.WarningSize, opts.CriticalAge, opts.CriticalSize)

	result := checkers.OK

	mtime := stat.ModTime()
	age := time.Now().Unix() - mtime.Unix()
	size := stat.Size()

	if monitor.CheckWarning(age, size) {
		result = checkers.WARNING
	}

	if monitor.CheckCritical(age, size) {
		result = checkers.CRITICAL
	}

	msg := fmt.Sprintf("%s is %d seconds old (%02d:%02d:%02d) and %d bytes.", opts.File, age, mtime.Hour(), mtime.Minute(), mtime.Second(), size)
	return checkers.NewChecker(result, msg)
}
