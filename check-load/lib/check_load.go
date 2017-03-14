package checkload

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var opts struct {
	WarningThreshold  string `short:"w" long:"warning" required:"true" value-name:"WL1,WL5,WL15" description:"Warning threshold for loadavg1,5,15"`
	CriticalThreshold string `short:"c" long:"critical" required:"true" value-name:"CL1,CL5,CL15" description:"Critical threshold for loadavg1,5,15"`
	PerCPU            bool   `short:"r" long:"percpu" description:"Divide the load averages by cpu count"`
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

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "LOAD"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	wload, err := parseThreshold(opts.WarningThreshold)
	if err != nil {
		return checkers.Unknown(err.Error())
	}
	cload, err := parseThreshold(opts.CriticalThreshold)
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	loadavgs, err := getloadavg()
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	result := checkers.OK
	for i, load := range loadavgs {
		if opts.PerCPU {
			numCPU := runtime.NumCPU()
			load = load / float64(numCPU)
		}
		if load > cload[i] {
			result = checkers.CRITICAL
			break
		}
		if load > wload[i] {
			result = checkers.WARNING
		}
	}

	msg := fmt.Sprintf("load average: %.2f, %.2f, %.2f", loadavgs[0], loadavgs[1], loadavgs[2])
	return checkers.NewChecker(result, msg)
}
