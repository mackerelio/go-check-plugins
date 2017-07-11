package checkdisk

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	Warning  *string `short:"w" long:"warning" value-name:"N, N%" description:"Exit with WARNING status if less than N units or N% of disk are free"`
	Critical *string `short:"c" long:"critical" value-name:"N, N%" description:"Exit with CRITICAL status if less than N units or N% of disk are free"`
	Path     *string `short:"p" long:"path" value-name:"PATH" description:"Mount point or block device as emitted by the mount(8) command"`
	Units    *string `short:"u" long:"units" value-name:"STRING" description:"Choose bytes, kB, MB, GB, TB (default: MB)"`
}

const (
	b  = float64(1)
	kb = float64(1024) * b
	mb = float64(1024) * kb
	gb = float64(1024) * mb
	tb = float64(1024) * gb
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

func checkStatus(val string, units float64, disk *diskStatus) (checkers.Status, error) {
	avail := float64(disk.Avail) / float64(units)
	freePct := (float64(disk.Avail) * float64(100)) / float64(disk.All)

	checkSt := checkers.OK
	if strings.HasSuffix(val, "%") {
		v, err := strconv.Atoi(strings.TrimRight(val, "%"))
		if err != nil {
			return checkers.UNKNOWN, err
		}

		if float64(v) > freePct {
			checkSt = checkers.WARNING
		}
	} else {
		v, err := strconv.Atoi(val)
		if err != nil {
			return checkers.UNKNOWN, err
		}

		if float64(v) > avail {
			checkSt = checkers.WARNING
		}
	}

	return checkSt, nil
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
		u := strings.ToLower(*opts.Units)
		if u == "bytes" {
			units = b
		} else if u == "kb" {
			units = kb
		} else if u == "gb" {
			units = gb
		} else if u == "tb" {
			units = tb
		} else {
			return checkers.Unknown(fmt.Sprintf("Faild to fetch disk usage: %s", errors.New("Invalid argument flag '-u, --units'")))
		}
	}

	checkSt := checkers.OK
	if opts.Warning != nil {
		checkSt, err = checkStatus(*opts.Warning, units, disk)
		if err != nil {
			return checkers.Unknown(fmt.Sprintf("Faild to check disk status: %s", err))
		}
	}

	if opts.Critical != nil {
		checkSt, err = checkStatus(*opts.Critical, units, disk)
		if err != nil {
			return checkers.Unknown(fmt.Sprintf("Faild to check disk status: %s", err))
		}
	}

	us := "MB"
	if opts.Units != nil {
		us = *opts.Units
	}

	all := float64(disk.All) / float64(units)
	used := float64(disk.Used) / float64(units)
	free := float64(disk.Free) / float64(units)
	avail := float64(disk.Avail) / float64(units)
	freePct := (float64(disk.Avail) * float64(100)) / float64(disk.All)
	msg := fmt.Sprintf("All: %.2f %v, Used: %.2f %v, Free: %.2f %v, Available: %.2f %v, Free percentage: %.2f\n", all, us, used, us, free, us, avail, us, freePct)

	return checkers.NewChecker(checkSt, msg)
}
