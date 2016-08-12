package main

import (
	"fmt"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var opts struct {
	ServiceName string `long:"service-name" short:"s" required:"true" description:"service name"`
	ListService bool   `long:"list-service" short:"l" description:"list service"`
}

type serviceState struct {
	Node             string `csv:"Node"`
	Caption          string `csv:"Caption"`
	DelayedAutoStart bool   `csv:"DelayedAutoStart"`
	ErrorControl     string `csv:"ErrorControl"`
	Name             string `csv:"Name"`
	Started          bool   `csv:"Started"`
	StartMode        string `csv:"StartMode"`
	State            string `csv:"State"`
}

func main() {
	ckr := run(os.Args[1:])
	ckr.Name = "NtService"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	ss, err := getServiceState()
	if opts.ListService {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if err != nil {
		return checkers.Critical(err.Error())
	}

	checkSt := checkers.OK
	msg := ""
	for _, s := range ss {
		if !strings.Contains(s.Name, opts.ServiceName) {
			continue
		}
		if s.State == "Running" {
			continue
		}
		if s.ErrorControl == "Ignore" {
			continue
		}
		checkSt = checkers.CRITICAL
		msg = fmt.Sprintf("%s: %s - %s", s.Name, s.Caption, s.State)
		break
	}

	return checkers.NewChecker(checkSt, msg)
}
