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
	gpud "github.com/shirou/gopsutil/disk"
)

type diskStatus struct {
	Dev   string
	All   uint64
	Used  uint64
	Free  uint64
	Avail uint64
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

type unit struct {
	Name string
	Size float64
}

func getDiskUsage(partition gpud.PartitionStat) (*diskStatus, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(partition.Mountpoint, &fs)
	if err != nil {
		return nil, err
	}

	disk := &diskStatus{}
	disk.Dev = partition.Device
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	disk.Avail = fs.Bavail * uint64(fs.Bsize)

	return disk, nil
}

func checkStatus(current checkers.Status, threshold string, units float64, disk *diskStatus, status checkers.Status) (checkers.Status, error) {
	avail := float64(disk.Avail) / float64(units)
	freePct := (float64(disk.Avail) * float64(100)) / float64(disk.All)

	if strings.HasSuffix(threshold, "%") {
		v, err := strconv.ParseFloat(strings.TrimRight(threshold, "%"), 64)
		if err != nil {
			return checkers.UNKNOWN, err
		}

		if v > freePct {
			current = status
		}
	} else {
		v, err := strconv.ParseFloat(threshold, 64)
		if err != nil {
			return checkers.UNKNOWN, err
		}

		if v > avail {
			current = status
		}
	}

	return current, nil
}

func genMessage(disk *diskStatus, u unit) string {
	all := float64(disk.All) / u.Size
	used := float64(disk.Used) / u.Size
	free := float64(disk.Free) / u.Size
	avail := float64(disk.Avail) / u.Size
	freePct := (float64(disk.Avail) * float64(100)) / float64(disk.All)

	return fmt.Sprintf("Dev: %v, All: %.2f %v, Used: %.2f %v, Free: %.2f %v, Available: %.2f %v, Free percentage: %.2f", disk.Dev, all, u.Name, used, u.Name, free, u.Name, avail, u.Name, freePct)
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

	partitions, err := gpud.Partitions(true)
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Faild to fetch partitions: %s", err))
	}

	if opts.Path != nil {
		contains := false
		for _, partition := range partitions {
			if *opts.Path == partition.Mountpoint {
				partitions = make([]gpud.PartitionStat, 0)
				partitions = append(partitions, partition)
				contains = true
			}

			if contains == false {
				return checkers.Unknown(fmt.Sprintf("Faild to fetch mountpoint: %s", errors.New("Invalid argument flag '-p, --path'")))
			}
		}
	}

	var disks []*diskStatus

	for _, partition := range partitions {
		disk, err := getDiskUsage(partition)
		if err != nil {
			return checkers.Unknown(fmt.Sprintf("Faild to fetch disk usage: %s", err))
		}

		disks = append(disks, disk)
	}

	u := unit{"MB", mb}
	if opts.Units != nil {
		us := strings.ToLower(*opts.Units)
		if us == "bytes" {
			u = unit{us, b}
		} else if us == "kb" {
			u = unit{us, mb}
		} else if us == "gb" {
			u = unit{us, gb}
		} else if us == "tb" {
			u = unit{us, tb}
		} else {
			return checkers.Unknown(fmt.Sprintf("Faild to check disk status: %s", errors.New("Invalid argument flag '-u, --units'")))
		}
	}

	checkSt := checkers.OK
	if opts.Warning != nil {
		for _, disk := range disks {
			checkSt, err = checkStatus(checkSt, *opts.Warning, u.Size, disk, checkers.WARNING)
			if err != nil {
				return checkers.Unknown(fmt.Sprintf("Faild to check disk status: %s", err))
			}

			if checkSt == checkers.WARNING {
				break
			}
		}
	}

	if opts.Critical != nil {
		for _, disk := range disks {
			checkSt, err = checkStatus(checkSt, *opts.Critical, u.Size, disk, checkers.CRITICAL)
			if err != nil {
				return checkers.Unknown(fmt.Sprintf("Faild to check disk status: %s", err))
			}

			if checkSt == checkers.CRITICAL {
				break
			}
		}
	}

	var msgs []string
	for _, disk := range disks {
		msg := genMessage(disk, u)
		msgs = append(msgs, msg)
	}
	msgss := strings.Join(msgs, ";\n")

	return checkers.NewChecker(checkSt, msgss)
}
