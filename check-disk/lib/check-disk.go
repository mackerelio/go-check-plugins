package checkdisk

import (
	"fmt"
	"os"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type diskStatus struct {
	All      uint64 `json:"all"`
	Used     uint64 `json:"used"`
	Free     uint64 `json:"free"`
	Avail    uint64 `json:"available"`
	FreeRate uint64 `json:"free_rate"`
}

var opts struct {
	Warning      *float64 `short:"W" long:"warning" value-name:"N" description:"Exit with WARNING status if less than N GB of disk are free"`
	Critical     *float64 `short:"C" long:"critical" value-name:"N" description:"Exit with CRITICAL status if less than N GB of disk are free"`
	WarningRate  *float64 `short:"w" long:"warning-rate" value-name:"N" description:"Exit with WARNING status if less than N % of disk are free"`
	CriticalRate *float64 `short:"c" long:"critical-rate" value-name:"N" description:"Exit with CRITICAL status if less than N % of disk are free"`
	Path         *string  `short:"p" long:"path" value-name:"PATH" description:"Mount point or block device as emitted by the mount(8) command"`
}

const (
	b  = 1
	kb = 1024 * b
	mb = 1024 * kb
	gb = 1024 * mb
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
	disk.FreeRate = (disk.Free / disk.All) * 100

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
		fmt.Printf("%v\n", err)
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

	all := float64(disk.All) / float64(gb)
	used := float64(disk.Used) / float64(gb)
	free := float64(disk.Free) / float64(gb)
	avail := float64(disk.Avail) / float64(gb)
	freeRate := float64(disk.FreeRate)

	checkSt := checkers.OK
	if opts.Warning != nil && *opts.Warning > free {
		checkSt = checkers.WARNING
	}
	if opts.Critical != nil && *opts.Critical > free {
		checkSt = checkers.CRITICAL
	}
	if opts.WarningRate != nil && *opts.WarningRate > freeRate {
		checkSt = checkers.WARNING
	}
	if opts.CriticalRate != nil && *opts.CriticalRate > freeRate {
		checkSt = checkers.CRITICAL
	}

	msg := fmt.Sprintf("All: %.2f GB, Used: %.2f GB, Free: %.2f GB, Available: %.2f GB, Free percentage: %.2f\n", all, used, free, avail, freeRate)

	return checkers.NewChecker(checkSt, msg)
}
