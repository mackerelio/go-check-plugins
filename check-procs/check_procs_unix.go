// +build !windows

package main

import (
	"errors"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var threadsUnknown = runtime.GOOS == "darwin"

func getProcs() (proc []procState, err error) {
	var procs []procState
	psformat := "user,pid,vsz,rss,pcpu,nlwp,state,etime,time,command"
	if threadsUnknown {
		psformat = "user,pid,vsz,rss,pcpu,state,etime,time,command"
	}
	output, _ := exec.Command("ps", "axwwo", psformat).Output()
	for _, line := range strings.Split(string(output), "\n")[1:] {
		proc, err := parseProcState(line)
		if err != nil {
			continue
		}
		procs = append(procs, proc)
	}
	return procs, nil
}

func parseProcState(line string) (proc procState, err error) {
	fields := strings.Fields(line)
	fieldsMinLen := 10
	if threadsUnknown {
		fieldsMinLen = 9
	}
	if len(fields) < fieldsMinLen {
		return procState{}, errors.New("parseProcState: insufficient words")
	}
	vsz, _ := strconv.ParseInt(fields[2], 10, 64)
	rss, _ := strconv.ParseInt(fields[3], 10, 64)
	pcpu, _ := strconv.ParseFloat(fields[4], 64)
	if threadsUnknown {
		esec := timeStrToSeconds(fields[6])
		csec := timeStrToSeconds(fields[7])
		return procState{strings.Join(fields[8:], " "), fields[0], fields[1], vsz, rss, pcpu, 1, fields[5], esec, csec}, nil
	}
	thcount, _ := strconv.ParseInt(fields[5], 10, 64)
	esec := timeStrToSeconds(fields[7])
	csec := timeStrToSeconds(fields[8])
	return procState{strings.Join(fields[9:], " "), fields[0], fields[1], vsz, rss, pcpu, thcount, fields[6], esec, csec}, nil
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

var timeRegexp = regexp.MustCompile(`(?:(\d+)-)?(?:(\d+):)?(\d+)[:.](\d+)`)

func timeStrToSeconds(etime string) int64 {
	match := timeRegexp.FindStringSubmatch(etime)
	if match == nil || len(match) != 5 {
		return 0
	}
	days, _ := strconv.ParseInt(match[1], 10, 64)
	hours, _ := strconv.ParseInt(match[2], 10, 64)
	minutes, _ := strconv.ParseInt(match[3], 10, 64)
	seconds, _ := strconv.ParseInt(match[4], 10, 64)
	return (((days*24+hours)*60+minutes)*60 + seconds)
}

