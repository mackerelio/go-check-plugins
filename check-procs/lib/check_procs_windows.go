package checkprocs

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

// Win32PerfFormattedDataPerfProcProcess is struct for Win32_PerfFormattedData_PerfProc_Process.
type Win32PerfFormattedDataPerfProcProcess struct {
	ElapsedTime          uint64
	IDProcess            uint32
	Name                 string
	PercentProcessorTime uint64
	ThreadCount          uint64
	VirtualBytes         uint64
	WorkingSet           uint64
}

func getProcs() (proc []procState, err error) {
	var records []Win32PerfFormattedDataPerfProcProcess
	err = wmi.Query("SELECT * FROM Win32_PerfFormattedData_PerfProc_Process WHERE Name != '_Total'", &records)
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		proc = append(proc, procState{
			cmd:     record.Name,
			pid:     fmt.Sprint(record.IDProcess),
			vsz:     int64(record.VirtualBytes),
			rss:     int64(record.WorkingSet),
			pcpu:    float64(record.PercentProcessorTime),
			thcount: int64(record.ThreadCount),
			esec:    int64(record.ElapsedTime),
			csec:    0,
		})
	}
	return proc, nil
}
