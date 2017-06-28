package checkdisk

import (
	"fmt"
	"os"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type diskStatus struct {
	All   uint64 `json:"all"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
	Avail uint64 `json:"available"`
}

var opts struct {
	Warning      *float64 `short:"W" long:"warning" value-name:"N" description:"Exit with WARNING status if less than N units of disk are free"`
	Critical     *float64 `short:"C" long:"critical" value-name:"N" description:"Exit with CRITICAL status if less than N units of disk are free"`
	WarningRate  *float64 `short:"w" long:"warning-rate" value-name:"N" description:"Exit with WARNING status if less than N % of disk are free"`
	CriticalRate *float64 `short:"c" long:"critical-rate" value-name:"N" description:"Exit with CRITICAL status if less than N % of disk are free"`
	Path         *string  `short:"p" long:"path" value-name:"PATH" description:"Mount point or block device as emitted by the mount(8) command"`
	Units        *string  `short:"u" long:"units" value-name:"STRING" description:"Choose bytes, kB, MB, GB, TB (default: MB)"`
}

const (
	b  = 1
	kb = 1024 * b
	mb = 1024 * kb
	gb = 1024 * mb
	tb = 1024 * gb
)

func getDiskUsage(path string) (*diskStatus, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return nil, err
	}

	disk := &diskStatus{}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	disk.Avail = fs.Bavail * uint64(fs.Bsize)

	return disk, nil
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "Disk"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	path := "/"
	if opts.Path != nil {
		path = *opts.Path
	}

	disk, err := getDiskUsage(path)
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Faild to fetch disk usage: %s", err))
	}

	units := mb
	if opts.Units != nil {
		switch {
		case *opts.Units == "bytes":
			units = b
		case *opts.Units == "kB":
			units = kb
		case *opts.Units == "MB":
			units = mb
		case *opts.Units == "GB":
			units = gb
		case *opts.Units == "TB":
			units = tb
		default:
			units = mb
		}
	}

	all := float64(disk.All) / float64(units)
	used := float64(disk.Used) / float64(units)
	free := float64(disk.Free) / float64(units)
	avail := float64(disk.Avail) / float64(units)
	freeRate := (float64(disk.Avail) * float64(100)) / float64(disk.All)

	checkSt := checkers.OK
	if opts.Warning != nil && *opts.Warning > avail {
		checkSt = checkers.WARNING
	}
	if opts.Critical != nil && *opts.Critical > avail {
		checkSt = checkers.CRITICAL
	}
	if opts.WarningRate != nil && *opts.WarningRate > freeRate {
		checkSt = checkers.WARNING
	}
	if opts.CriticalRate != nil && *opts.CriticalRate > freeRate {
		checkSt = checkers.CRITICAL
	}

	u := "MB"
	if opts.Units != nil {
		u = *opts.Units
	}

	msg := fmt.Sprintf("All: %.2f %v, Used: %.2f %v, Free: %.2f %v, Available: %.2f %v, Free percentage: %.2f\n", all, u, used, u, free, u, avail, u, freeRate)

	return checkers.NewChecker(checkSt, msg)
}
