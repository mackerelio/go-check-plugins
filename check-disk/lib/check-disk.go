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
	Warning  *float64 `short:"w" long:"warning" value-name:"N" description:"Exit with WARNING status if less than INTEGER GB of disk are free"`
	Critical *float64 `short:"c" long:"critical" value-name:"N" description:"Exit with CRITICAL status if less than INTEGER GB of disk are free"`
	Path     *string  `short:"p" long:"path" value-name:"PATH" description:"Mount point or block device as emitted by the mount(8) command"`
}

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

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

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

	checkSt := checkers.OK
	if opts.Warning != nil && *opts.Warning > float64(disk.Free)/float64(GB) {
		checkSt = checkers.WARNING
	}
	if opts.Critical != nil && *opts.Critical > float64(disk.Free)/float64(GB) {
		checkSt = checkers.CRITICAL
	}

	all := float64(disk.All) / float64(GB)
	used := float64(disk.Used) / float64(GB)
	free := float64(disk.Free) / float64(GB)
	avail := float64(disk.Avail) / float64(GB)

	msg := fmt.Sprintf("All: %.2f GB, Used: %.2f GB, Free: %.2f GB, Available: %.2f GB\n", all, used, free, avail)

	return checkers.NewChecker(checkSt, msg)
}
