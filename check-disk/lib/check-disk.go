package checkdisk

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	gpud "github.com/shirou/gopsutil/disk"
)

var opts struct {
	Warning  *string `short:"w" long:"warning" value-name:"N, N%" description:"Exit with WARNING status if less than N units or N% of disk are free"`
	Critical *string `short:"c" long:"critical" value-name:"N, N%" description:"Exit with CRITICAL status if less than N units or N% of disk are free"`
	Path     *string `short:"p" long:"path" value-name:"PATH" description:"Mount point or block device as emitted by the mount(8) command"`
	Exclude  *string `short:"x" long:"exclude_device" value-name:"EXCLUDE PATH" description:"Ignore device (only works if -p unspecified)"`
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

func checkStatus(current checkers.Status, threshold string, units float64, disk *gpud.UsageStat, status checkers.Status) (checkers.Status, error) {
	if strings.HasSuffix(threshold, "%") {
		v, err := strconv.ParseFloat(strings.TrimRight(threshold, "%"), 64)
		if err != nil {
			return checkers.UNKNOWN, err
		}

		freePct := float64(100) - disk.UsedPercent
		inodesFreePct := float64(100) - disk.InodesUsedPercent

		if v > freePct || v > inodesFreePct {
			current = status
		}
	} else {
		v, err := strconv.ParseFloat(threshold, 64)
		if err != nil {
			return checkers.UNKNOWN, err
		}

		if v > float64(disk.Free) {
			current = status
		}
	}

	return current, nil
}

func genMessage(disk *gpud.UsageStat, u unit) string {
	all := float64(disk.Total) / u.Size
	used := float64(disk.Used) / u.Size
	free := float64(disk.Free) / u.Size
	freePct := float64(100) - disk.UsedPercent
	inodesFreePct := float64(100) - disk.InodesUsedPercent

	return fmt.Sprintf("Path: %v, All: %.2f %v, Used: %.2f %v, Free: %.2f %v, Free percentage: %.2f (inodes: %.2f)", disk.Path, all, u.Name, used, u.Name, free, u.Name, freePct, inodesFreePct)
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

	partitions, err := listPartitions()
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("Faild to fetch partitions: %s", err))
	}

	if opts.Path != nil {
		exist := false
		for _, partition := range partitions {
			if *opts.Path == partition.Mountpoint {
				partitions = make([]gpud.PartitionStat, 0)
				partitions = append(partitions, partition)
				exist = true
			}
		}

		if exist == false {
			return checkers.Unknown(fmt.Sprintf("Faild to fetch mountpoint: %s", errors.New("Invalid argument flag '-p, --path'")))
		}
	}

	if opts.Path == nil && opts.Exclude != nil {
		var tmp []gpud.PartitionStat
		for _, partition := range partitions {
			if *opts.Exclude != partition.Mountpoint {
				tmp = append(tmp, partition)
			}
		}
		partitions = tmp
	}

	var disks []*gpud.UsageStat

	for _, partition := range partitions {
		disk, err := gpud.Usage(partition.Mountpoint)
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

// ref: mountlist.c in gnulib
// https://github.com/coreutils/gnulib/blob/a742bdb3/lib/mountlist.c#L168
func listPartitions() ([]gpud.PartitionStat, error) {
	allPartitions, err := gpud.Partitions(true)
	if err != nil {
		return nil, err
	}
	partitions := make([]gpud.PartitionStat, 0, len(allPartitions))
	for _, p := range allPartitions {
		switch p.Fstype {
		case "autofs":
		case "proc":
		case "subfs":
		case "debugfs":
		case "devpts":
		case "fusectl":
		case "mqueue":
		case "rpc_pipefs":
		case "sysfs":
		case "devfs":
		case "kernfs":
		case "ignore":
			continue
		case "none":
			if !strings.Contains(p.Opts, "bind") {
				partitions = append(partitions, p)
			}
		default:
			partitions = append(partitions, p)
		}
	}
	return partitions, nil
}
