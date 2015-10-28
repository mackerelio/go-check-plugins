package main

import (
	"encoding/csv"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

func getProcs() (proc []procState, err error) {
	var procs []procState
	// WMIC PATH Win32_PerfFormattedData_PerfProc_Process WHERE "Name != '_Total'" GET Name,IDProcess,VirtualBytes,WorkingSet,PercentProcessorTime,ThreadCount,ElapsedTime /FORMAT:CSV
	output, _ := exec.Command("WMIC", "PATH", "Win32_PerfFormattedData_PerfProc_Process", "WHERE", "Name != '_Total'", "GET", "ElapsedTime,IDProcess,Name,PercentProcessorTime,ThreadCount,VirtualBytes,WorkingSet", "/FORMAT:CSV").Output()
	r := csv.NewReader(strings.NewReader(string(output[1:])))
	records, err := r.ReadAll()
	if (err != nil) {
		return procs, nil
	}
	for _, record := range records[1:] {
		proc, err := parsePerfProc(record)
		if err != nil {
			continue
		}
		procs = append(procs, proc)
	}
	return procs, nil
}

func parsePerfProc(fields []string) (proc procState, err error) {
	fieldsLen := 8
	if len(fields) != fieldsLen {
		return procState{}, errors.New("parseTaskList: insufficient words")
	}
	vsz, _ := strconv.ParseInt(fields[6], 10, 64) //VirtualBytes
	rss, _ := strconv.ParseInt(fields[7], 10, 64) // WorkingSet
	pcpu, _ := strconv.ParseFloat(fields[4], 64) // PercentProcessorTime
	thcount, _ := strconv.ParseInt(fields[5], 10, 64) //ThreadCount
	esec, _ := strconv.ParseInt(fields[1], 10, 64) // ElapsedTime
	csec := int64(0)
	return procState{fields[3] /* Name */, "", fields[2] /* IDProcess */, vsz, rss, pcpu, thcount, "", esec, csec}, nil
}
