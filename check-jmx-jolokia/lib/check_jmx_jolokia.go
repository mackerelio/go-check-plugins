package checkjmxjolokia

import (
	"encoding/json"
	"fmt"
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

type jmxJolokiaResponse struct {
	Status int
	Value  float64
}

// Do the plugin
func Do() {
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

	resJ := jmxJolokiaResponse{}
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&resJ); err != nil {
		return checkers.Critical(err.Error())
	}

	if resJ.Status != 200 {
		return checkers.Unknown(fmt.Sprintf("failed: response status %d", resJ.Status))
	}

	checkSt := checkers.OK
	msg := fmt.Sprintf("%s %s value %f", opts.MBean, opts.Attribute, resJ.Value)
	if resJ.Value > opts.Critical {
		checkSt = checkers.CRITICAL
		msg = fmt.Sprintf("%s %s value is over %f > %f", opts.MBean, opts.Attribute, resJ.Value, opts.Critical)
	} else if resJ.Value > opts.Warning {
		checkSt = checkers.WARNING
		msg = fmt.Sprintf("%s %s value is over %f > %f", opts.MBean, opts.Attribute, resJ.Value, opts.Warning)
	}

	return checkers.NewChecker(checkSt, msg)
}