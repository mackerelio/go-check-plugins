package checkuptime

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/mackerelio/go-osstat/uptime"
	"github.com/mackerelio/golib/pluginutil"
)

type uptimeOpts struct {
	WarnUnder     *float64 `long:"warn-under" value-name:"N" description:"(DEPRECATED) Trigger a warning if under the seconds"`
	WarningUnder  *float64 `short:"w" long:"warning-under" value-name:"N" description:"Trigger a warning if under the seconds"`
	CritUnder     *float64 `short:"c" long:"critical-under" value-name:"N" description:"Trigger a critial if under the seconds"`
	WarnOver      *float64 `long:"warn-over" value-name:"N" description:"(DEPRECATED) Trigger a warning if over the seconds"`
	WarningOver   *float64 `short:"W" long:"warning-over" value-name:"N" description:"Trigger a warning if over the seconds"`
	CritOver      *float64 `short:"C" long:"critical-over" value-name:"N" description:"Trigger a critical if over the seconds"`
	WarningRewind bool     `short:"r" long:"warning-rewind" description:"Trigger a warning if rewind uptime (detect reboot)"`
	CritRewind    bool     `short:"R" long:"critical-rewind" description:"Trigger a critical if rewind uptime (detect reboot)"`
	StateDir      string   `short:"s" long:"state-dir" value-name:"DIR" description:"Dir to keep state files under"`
	origArgs      []string
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "Uptime"
	ckr.Exit()
}

func parseArgs(args []string) (*uptimeOpts, error) {
	origArgs := make([]string, len(args))
	copy(origArgs, args)
	opts := &uptimeOpts{}
	_, err := flags.ParseArgs(opts, args)
	opts.origArgs = origArgs
	if opts.StateDir == "" {
		workdir := pluginutil.PluginWorkDir()
		opts.StateDir = filepath.Join(workdir, "check-uptime")
	}
	return opts, err
}

func getStateFile(stateDir string, args []string) string {
	hash := md5.Sum([]byte(strings.Join(args, " ")))
	return filepath.Join(
		stateDir,
		hex.EncodeToString(hash[:]),
	)
}

func getPreviousUptime(f string) (float64, error) {
	_, err := os.Stat(f)
	if err != nil {
		return 0, nil
	}
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return 0, err
	}
	pUt, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return 0, err
	}
	return pUt, nil
}

func writePreviousUptime(f string, pUt float64) error {
	err := os.MkdirAll(filepath.Dir(f), 0755)
	if err != nil {
		return err
	}
	return writeFileAtomically(f, []byte(strconv.FormatFloat(pUt, 'f', -1, 64)))
}

func writeFileAtomically(f string, contents []byte) error {
	// MUST be located on same disk partition
	tmpf, err := ioutil.TempFile(filepath.Dir(f), "tmp")
	if err != nil {
		return err
	}
	// os.Remove here works successfully when tmpf.Write fails or os.Rename fails.
	// In successful case, os.Remove fails because the temporary file is already renamed.
	defer os.Remove(tmpf.Name())
	_, err = tmpf.Write(contents)
	tmpf.Close() // should be called before rename
	if err != nil {
		return err
	}
	return os.Rename(tmpf.Name(), f)
}

func run(args []string) *checkers.Checker {
	opts, err := parseArgs(args)
	if err != nil {
		os.Exit(1)
	}
	utDur, err := uptime.Get()
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Faild to fetch uptime metrics: %s", err))
	}
	ut := utDur.Seconds()

	// for backward compatibility
	if opts.WarnUnder != nil && opts.WarningUnder == nil {
		opts.WarningUnder = opts.WarnUnder
	}
	if opts.WarnOver != nil && opts.WarningOver == nil {
		opts.WarningOver = opts.WarnOver
	}

	checkSt := checkers.OK
	if opts.WarningRewind || opts.CritRewind {
		stateFile := getStateFile(opts.StateDir, opts.origArgs)
		pUt, err := getPreviousUptime(stateFile)
		if err != nil {
			return checkers.Unknown(fmt.Sprintf("Failed to get previous uptime: %s", err))
		}

		if opts.WarningRewind && pUt > ut {
			checkSt = checkers.WARNING
		}
		if opts.CritRewind && pUt > ut {
			checkSt = checkers.CRITICAL
		}
		writePreviousUptime(stateFile, ut)
	}
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
