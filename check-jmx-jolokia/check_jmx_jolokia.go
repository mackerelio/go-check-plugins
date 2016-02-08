package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type jmxJolokiaOpts struct {
	HostName  string  `short:"H" long:"host" required:"true" description:"Host name or IP Address"`
	Port      int     `short:"p" long:"port" default:"8778" description:"Port"`
	Timeout   int     `short:"t" long:"timeout" default:"10" description:"Seconds before connection times out"`
	MBean     string  `short:"m" long:"mbean" required:"true" description:"MBean"`
	Attribute string  `short:"a" long:"attribute" required:"true" description:"Attribute"`
	InnerPath string  `short:"i" long:"inner-path" description:"InnerPath"`
	Key       string  `short:"k" long:"key" default:"value" description:"Key"`
	Warning   float64 `short:"w" long:"warning" description:"Trigger a warning if over a number"`
	Critical  float64 `short:"c" long:"critical" description:"Trigger a critical if over a number"`
}

type jmxJolokiaResult struct {
	Status int
	Value  float64
}

func main() {
	ckr := run(os.Args[1:])
	ckr.Name = "Jmx-Jolokia"
	ckr.Exit()
}

func parseArgs(args []string) (*jmxJolokiaOpts, error) {
	opts := &jmxJolokiaOpts{}
	_, err := flags.ParseArgs(opts, args)
	return opts, err
}

func createURL(opts *jmxJolokiaOpts) string {
	if opts.InnerPath == "" {
		return fmt.Sprintf("http://%s:%d/jolokia/read/%s/%s", opts.HostName, opts.Port, opts.MBean, opts.Attribute)
	}
	return fmt.Sprintf("http://%s:%d/jolokia/read/%s/%s/%s", opts.HostName, opts.Port, opts.MBean, opts.Attribute, opts.InnerPath)
}

func run(args []string) *checkers.Checker {
	opts, err := parseArgs(args)
	if err != nil {
		os.Exit(1)
	}

	client := &http.Client{Timeout: time.Duration(opts.Timeout) * time.Second}
	res, err := client.Get(createURL(opts))
	if err != nil {
		return checkers.Critical(err.Error())
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return checkers.Unknown(fmt.Sprintf("failed: http status code %d", res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	result := jmxJolokiaResult{}
	if err := json.Unmarshal(body, &result); err != nil {
		return checkers.Unknown(err.Error())
	}

	if result.Status != 200 {
		return checkers.Unknown(fmt.Sprintf("failed: response status %d", result.Status))
	}

	checkSt := checkers.OK
	msg := fmt.Sprintf("%s %s value %f", opts.MBean, opts.Attribute, result.Value)
	if result.Value > opts.Critical {
		checkSt = checkers.CRITICAL
		msg = fmt.Sprintf("%s %s value is over %f > %f", opts.MBean, opts.Attribute, result.Value, opts.Critical)
	} else if result.Value > opts.Warning {
		checkSt = checkers.WARNING
		msg = fmt.Sprintf("%s %s value is over %f > %f", opts.MBean, opts.Attribute, result.Value, opts.Warning)
	}

	return checkers.NewChecker(checkSt, msg)
}