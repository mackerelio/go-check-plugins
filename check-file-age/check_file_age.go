package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type checkStatus int

const (
	ok checkStatus = iota
	warning
	critical
	unknown
)

type monitor struct {
	warningAge   int64
	warningSize  int64
	criticalAge  int64
	criticalSize int64
}

func (m monitor) hasWarningAge() bool {
	return m.warningAge != 0
}

func (m monitor) hasWarningSize() bool {
	return m.warningSize != 0
}

func (m monitor) CheckWarning(age, size int64) bool {
	return (m.hasWarningAge() && m.warningAge < age) ||
		(m.hasWarningSize() && m.warningSize > size)
}

func (m monitor) hasCriticalAge() bool {
	return m.criticalAge != 0
}

func (m monitor) hasCriticalSize() bool {
	return m.criticalSize != 0
}

func (m monitor) CheckCritical(age, size int64) bool {
	return (m.hasCriticalAge() && m.criticalAge < age) ||
		(m.hasCriticalSize() && m.criticalSize > size)
}

func newMonitor(warningAge, warningSize, criticalAge, criticalSize int64) *monitor {
	return &monitor{
		warningAge:   warningAge,
		warningSize:  warningSize,
		criticalAge:  criticalAge,
		criticalSize: criticalSize,
	}
}

func main() {
	var (
		file          string
		warningAge    int64
		warningSize   int64
		criticalAge   int64
		criticalSize  int64
		ignoreMissing bool
	)

	flag.StringVar(&file, "f", "", "file")
	flag.StringVar(&file, "file", "", "file")
	flag.Int64Var(&warningAge, "w", 240, "warning age")
	flag.Int64Var(&warningAge, "warning-age", 240, "warning age")
	flag.Int64Var(&warningSize, "W", 0, "warning size")
	flag.Int64Var(&warningSize, "warning-size", 0, "warning size")
	flag.Int64Var(&criticalAge, "c", 600, "critical age")
	flag.Int64Var(&criticalAge, "critical-age", 600, "critical age")
	flag.Int64Var(&criticalSize, "C", 0, "critical size")
	flag.Int64Var(&criticalSize, "critical-size", 0, "critical size")
	flag.BoolVar(&ignoreMissing, "i", false, "ignore missing")
	flag.BoolVar(&ignoreMissing, "ignore-missing", false, "ignore missing")

	flag.Parse()

	if file == "" {
		if file = flag.Arg(0); file == "" {
			fmt.Println("No file specified")
			os.Exit(int(unknown))
		}
	}

	stat, err := os.Stat(file)
	if err != nil {
		if ignoreMissing {
			fmt.Println("No such file, but ignore missing is set.")
			os.Exit(int(ok))
		} else {
			fmt.Println(err.Error())
			os.Exit(int(unknown))
		}
	}

	monitor := newMonitor(warningAge, warningSize, criticalAge, criticalSize)

	result := ok

	age := time.Now().Unix() - stat.ModTime().Unix()
	size := stat.Size()

	if monitor.CheckWarning(age, size) {
		result = warning
	}

	if monitor.CheckCritical(age, size) {
		result = critical
	}

	fmt.Printf("%s is %d seconds old and %d bytes.\n", file, age, size)
	os.Exit(int(result))
}
