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

func main() {
	var (
		file          = flag.String("f", "", "file")
		warningAge    = flag.Int64("w", 240, "warning age")
		warningSize   = flag.Int64("W", 0, "warning size")
		criticalAge   = flag.Int64("c", 600, "critical age")
		criticalSize  = flag.Int64("C", 0, "critical size")
		ignoreMissing = flag.Bool("i", false, "ignore missing")
	)

	flag.Parse()

	if *file == "" {
		if *file = flag.Arg(0); *file == "" {
			fmt.Println("No file specified")
			os.Exit(int(unknown))
		}
	}

	stat, err := os.Stat(*file)
	if err != nil {
		if *ignoreMissing {
			fmt.Println("No such file, but ignore missing is set.")
			os.Exit(int(ok))
		} else {
			fmt.Println(err.Error())
			os.Exit(int(unknown))
		}
	}

	result := ok

	age := time.Now().Unix() - stat.ModTime().Unix()
	size := stat.Size()

	if (*warningAge != 0 && *warningAge <= age) || (*warningSize != 0 && *warningSize <= size) {
		result = warning
	}

	if (*criticalAge != 0 && *criticalAge <= age) || (*criticalSize != 0 && *criticalSize <= size) {
		result = critical
	}

	fmt.Printf("%s is %d seconds old and %d bytes.\n", *file, age, size)
	os.Exit(int(result))
}
