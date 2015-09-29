package main

// // for getloadavg(2)
/*
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	WarningThreshold  string `short:"w" long:"warning" required:"true" value-name:"WL1,WL5,WL15" description:"Warning threshold for loadavg1,5,15"`
	CriticalThreshold string `short:"c" long:"critical" required:"true" value-name:"CL1,CL5,CL15" description:"Critical threshold for loadavg1,5,15"`
	PerCPU            bool   `short:"r" long:"percpu" default:"false" description:"Divide the load averages by cpu count"`
}

type checkStatus int

const (
	ok checkStatus = iota
	warning
	critical
	unknown
)

func (chSt checkStatus) String() string {
	switch {
	case chSt == 0:
		return "OK"
	case chSt == 1:
		return "WARNING"
	case chSt == 2:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

func parseThreshold(str string) ([3]float64, error) {
	thresholds := [3]float64{0, 0, 0}

	thSt := strings.Split(str, ",")
	if len(thSt) != 3 {
		return thresholds, errors.New("Threshold must be comma-separated 3 numbers")
	}

	var err error
	for i, v := range thSt {
		thresholds[i], err = strconv.ParseFloat(v, 64)
		if err != nil {
			return thresholds, err
		}
	}
	return thresholds, nil
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Print(err)
		os.Exit(int(unknown))
	}

	wload, err := parseThreshold(opts.WarningThreshold)
	if err != nil {
		fmt.Print(err)
		os.Exit(int(unknown))
	}
	cload, err := parseThreshold(opts.CriticalThreshold)
	if err != nil {
		fmt.Print(err)
		os.Exit(int(unknown))
	}

	loadavgs, _ := getloadavg()

	result := ok
	for i, load := range loadavgs {
		if opts.PerCPU {
			numCPU := runtime.NumCPU()
			load = load / float64(numCPU)
		}
		if load > cload[i] {
			result = critical
			break
		}
		if load > wload[i] {
			result = warning
		}
	}

	fmt.Printf("%s - load average: %.2f, %.2f, %.2f\n", result.String(), loadavgs[0], loadavgs[1], loadavgs[2])
	os.Exit(int(result))
}
