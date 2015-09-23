package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var opts struct {
	URL string `short:"u" long:"url" description:"A URL to connect to"`
}

func main() {
	ckr := run(os.Args[1:])
	ckr.Name = "HTTP"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil || opts.URL == "" {
		os.Exit(1)
	}

	// XXX set UA and etc... check_http (go-check-plugins)
	stTime := time.Now()
	resp, err := http.Get(opts.URL)
	if err != nil {
		return checkers.Critical(err.Error())
	}
	elapsed := time.Since(stTime)
	defer resp.Body.Close()

	cLength := resp.ContentLength
	if cLength == -1 {
		byt, _ := ioutil.ReadAll(resp.Body)
		cLength = int64(len(byt))
	}

	checkSt := checkers.UNKNOWN
	switch st := resp.StatusCode; true {
	case st < 400:
		checkSt = checkers.OK
	case st < 500:
		checkSt = checkers.WARNING
	default:
		checkSt = checkers.CRITICAL
	}

	msg := fmt.Sprintf("%s %s - %d bytes in %f second respons time",
		resp.Proto, resp.Status, cLength, elapsed.Seconds())

	return checkers.NewChecker(checkSt, msg)
}
