package checkdisk

import (
	"testing"

	gpud "github.com/shirou/gopsutil/v3/disk"
)

var root = gpud.UsageStat{
	Path:              "/",
	Total:             6000,
	Free:              3000,
	Used:              3000,
	UsedPercent:       50.1,
	InodesTotal:       4000,
	InodesUsed:        5000,
	InodesFree:        6000,
	InodesUsedPercent: 49.1,
	Fstype:            "xfs",
}

var opt = gpud.UsageStat{
	Path:              "/opt",
	Total:             60000,
	Free:              5000,
	Used:              40000,
	UsedPercent:       91.7,
	InodesTotal:       60000,
	InodesUsed:        20000,
	InodesFree:        40000,
	InodesUsedPercent: 91.7,
	Fstype:            "xfs",
}

var data = gpud.UsageStat{
	Path:              "/pgdata",
	Total:             600000,
	Free:              10000,
	Used:              590000,
	UsedPercent:       98.3,
	InodesTotal:       600000,
	InodesUsed:        10000,
	InodesFree:        590000,
	InodesUsedPercent: 98.3,
	Fstype:            "xfs",
}

func TestSortDisks(t *testing.T) {
	disks := []*gpud.UsageStat{&root, &opt, &data}

	sortDisks(disks)

	if disks[0].Path != "/pgdata" {
		t.Errorf("Expected '/pgdata' as first element but received '%s'", disks[0].Path)
	}

	if disks[1].Path != "/opt" {
		t.Errorf("Expected 'opt' as second element but received '%s'", disks[1].Path)
	}
}
