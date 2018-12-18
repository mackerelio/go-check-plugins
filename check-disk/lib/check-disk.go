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
	Warning       *string   `short:"w" long:"warning" value-name:"N, N%" description:"Exit with WARNING status if less than N units or N% of disk are free"`
	Critical      *string   `short:"c" long:"critical" value-name:"N, N%" description:"Exit with CRITICAL status if less than N units or N% of disk are free"`
	InodeWarning  *string   `short:"W" long:"iwarning" value-name:"N%" description:"Exit with WARNING status if less than PERCENT of inode space is free"`
	InodeCritical *string   `short:"K" long:"icritical" value-name:"N%" description:"Exit with CRITICAL status if less than PERCENT of inode space is free"`
	Path          *[]string `short:"p" long:"path" value-name:"PATH" description:"Mount point or block device as emitted by the mount(8) command (may be repeated)"`
	Exclude       *[]string `short:"x" long:"exclude-device" value-name:"EXCLUDE PATH" description:"Ignore device (may be repeated; only works if -p unspecified)"`
	All           bool      `short:"A" long:"all" description:"Explicitly select all paths."`
	ExcludeType   *[]string `short:"X" long:"exclude-type" value-name:"TYPE" description:"Ignore all filesystems of indicated type (may be repeated)"`
	IncludeType   *[]string `short:"N" long:"include-type" value-name:"TYPE" description:"Check only filesystems of indicated type (may be repeated)"`
	Units         *string   `short:"u" long:"units" value-name:"STRING" description:"Choose bytes, kB, MB, GB, TB (default: MB)"`
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

func checkDisk(current checkers.Status, threshold string, units float64, disk *gpud.UsageStat, status checkers.Status) (checkers.Status, error) {
	if strings.HasSuffix(threshold, "%") {
		thresholdPct, err := strconv.ParseFloat(strings.TrimRight(threshold, "%"), 64)
		if err != nil {
			return checkers.UNKNOWN, err
		}

		freePct := float64(100) - disk.UsedPercent
		if thresholdPct > freePct {
			current = status
		}
	} else {
		thresholdVal, err := strconv.ParseFloat(threshold, 64)
		if err != nil {
			return checkers.UNKNOWN, err
		}

		free := float64(disk.Free) / units
		if thresholdVal > free {
			current = status
		}
	}

	return current, nil
}

func checkInodes(current checkers.Status, threshold string, disk *gpud.UsageStat, status checkers.Status) (checkers.Status, error) {
	if !strings.HasSuffix(threshold, "%") {
		return checkers.UNKNOWN, errors.New("-W, -K value should be N%")
	}

	thresholdPct, err := strconv.ParseFloat(strings.TrimRight(threshold, "%"), 64)
	if err != nil {
		return checkers.UNKNOWN, err
	}

	if disk.InodesTotal == 0 {
		return checkers.UNKNOWN, fmt.Errorf("Disk %s does not have inodes", disk.Path)
	}

	inodesFreePct := float64(100) - disk.InodesUsedPercent
	if thresholdPct > inodesFreePct {
		current = status
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
		return checkers.Unknown(fmt.Sprintf("Failed to fetch partitions: %s", err))
	}

	if !opts.All {
		// Filtering partitions by Fstype
		if opts.IncludeType != nil {
			partitions = filterPartitionsByInclusion(partitions, *opts.IncludeType, fstypeOfPartition)
			if len(partitions) == 0 {
				return checkers.Unknown(fmt.Sprintf("Failed to fetch partitions: %s", errors.New("No device found for the specified *FsType*")))
			}
		}

		if opts.ExcludeType != nil {
			partitions = filterPartitionsByExclusion(partitions, *opts.ExcludeType, fstypeOfPartition)
		}

		// Filtering partions by Mountpoint
		if opts.Path != nil {
			if opts.Exclude != nil {
				return checkers.Unknown(fmt.Sprintf("Invalid arguments: %s", errors.New("-x does not work with -p")))
			}

			partitions = filterPartitionsByInclusion(partitions, *opts.Path, mountpointOfPartition)
			if len(partitions) == 0 {
				return checkers.Unknown(fmt.Sprintf("Failed to fetch partitions: %s", errors.New("No device found for the specified *Mountpoint*")))
			}
		}

		if opts.Path == nil && opts.Exclude != nil {
			partitions = filterPartitionsByExclusion(partitions, *opts.Exclude, mountpointOfPartition)
		}
	}

	if len(partitions) == 0 {
		return checkers.Unknown(fmt.Sprintf("Failed to fetch partitions: %s", errors.New("No device found")))
	}

	var disks []*gpud.UsageStat

	for _, partition := range partitions {
		disk, err := gpud.Usage(partition.Mountpoint)
		if err != nil {
			return checkers.Unknown(fmt.Sprintf("Failed to fetch disk usage: %s", err))
		}

		if disk.Total != 0 {
			disks = append(disks, disk)
		}
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
			return checkers.Unknown(fmt.Sprintf("Failed to check disk status: %s", errors.New("Invalid argument flag '-u, --units'")))
		}
	}

	checkSt := checkers.OK
	if opts.InodeCritical != nil {
		for _, disk := range disks {
			checkSt, err = checkInodes(checkSt, *opts.InodeCritical, disk, checkers.CRITICAL)
			if err != nil {
				return checkers.Unknown(fmt.Sprintf("Failed to check disk status: %s", err))
			}

			if checkSt == checkers.CRITICAL {
				break
			}
		}
	}

	if checkSt != checkers.CRITICAL && opts.Critical != nil {
		for _, disk := range disks {
			checkSt, err = checkDisk(checkSt, *opts.Critical, u.Size, disk, checkers.CRITICAL)
			if err != nil {
				return checkers.Unknown(fmt.Sprintf("Failed to check disk status: %s", err))
			}

			if checkSt == checkers.CRITICAL {
				break
			}
		}
	}

	if checkSt != checkers.CRITICAL && opts.InodeWarning != nil {
		for _, disk := range disks {
			checkSt, err = checkInodes(checkSt, *opts.InodeWarning, disk, checkers.WARNING)
			if err != nil {
				return checkers.Unknown(fmt.Sprintf("Failed to check disk status: %s", err))
			}

			if checkSt == checkers.WARNING {
				break
			}
		}
	}

	if checkSt == checkers.OK && opts.Warning != nil {
		for _, disk := range disks {
			checkSt, err = checkDisk(checkSt, *opts.Warning, u.Size, disk, checkers.WARNING)
			if err != nil {
				return checkers.Unknown(fmt.Sprintf("Failed to check disk status: %s", err))
			}

			if checkSt == checkers.WARNING {
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
		case "autofs",
			"proc",
			"subfs",
			"debugfs",
			"devpts",
			"fusectl",
			"mqueue",
			"rpc_pipefs",
			"sysfs",
			"devfs",
			"kernfs",
			"ignore":
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

func mountpointOfPartition(partition gpud.PartitionStat) string {
	return partition.Mountpoint
}

func fstypeOfPartition(partition gpud.PartitionStat) string {
	return partition.Fstype
}

func filterPartitionsByInclusion(partitions []gpud.PartitionStat, list []string, key func(_ gpud.PartitionStat) string) []gpud.PartitionStat {
	newPartitions := make([]gpud.PartitionStat, 0, len(partitions))
	for _, partition := range partitions {
		var ok = false
		for _, l := range list {
			if l == key(partition) {
				ok = true
				break
			}
		}
		if ok {
			newPartitions = append(newPartitions, partition)
		}
	}

	return newPartitions
}

func filterPartitionsByExclusion(partitions []gpud.PartitionStat, list []string, key func(_ gpud.PartitionStat) string) []gpud.PartitionStat {
	newPartitions := make([]gpud.PartitionStat, 0, len(partitions))
	for _, partition := range partitions {
		var ok = true
		for _, l := range list {
			if l == key(partition) {
				ok = false
				break
			}
		}
		if ok {
			newPartitions = append(newPartitions, partition)
		}
	}

	return newPartitions
}
