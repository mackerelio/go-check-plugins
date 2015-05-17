package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	URL string `short:"u" long:"url" description:"A URL to connect to"`
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

func main() {
	_, err := flags.Parse(&opts)
	if err != nil || opts.URL == "" {
		os.Exit(1)
	}

	// XXX set UA and etc... check_http (go-check-plugins)
	stTime := time.Now()
	resp, err := http.Get(opts.URL)
	if err != nil {
		fmt.Printf("HTTP CRITICAL %s\n", err)
		os.Exit(2)
	}
	elapsed := time.Since(stTime)

	defer resp.Body.Close()
	checkSt := unknown
	switch st := resp.StatusCode; true {
	case st < 400:
		checkSt = ok
	case st < 500:
		checkSt = warning
	default:
		checkSt = critical
	}

	fmt.Printf("HTTP %s: %s %s - %d bytes in %f second respons time",
		checkSt.String(), resp.Proto, resp.Status, resp.ContentLength, elapsed.Seconds())
	os.Exit(int(checkSt))
}
