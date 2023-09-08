package checkntservice

import (
	"fmt"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var opts struct {
	ServiceName    string `long:"service-name" short:"s" description:"matches if contained in service name."`
	ExcludeService string `long:"exclude-service" short:"x" description:"exclude if contained in service name. This option takes precedence over --service-name."`
	ListService    bool   `long:"list-service" short:"l" description:"list service"`
	Exact          bool   `long:"exact" description:"more exact checking of the service. This option applies only to --service-name."`
}

// Win32Service is struct for Win32_Service.
type Win32Service struct {
	Caption string
	Name    string
	State   string
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "NtService"
	ckr.Exit()
}

var getServiceStateFunc = getServiceState

func run(args []string) *checkers.Checker {
	var parser = flags.NewParser(&opts, flags.Default)
	_, err := parser.ParseArgs(args)
	if err != nil {
		os.Exit(1)
	}

	ss, err := getServiceStateFunc()
	if opts.ListService {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, s := range ss {
			fmt.Printf("%s: %s\n", s.Name, s.Caption)
		}
		os.Exit(0)
	}
	if opts.ServiceName == "" {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	if err != nil {
		return checkers.Critical(err.Error())
	}

	checkSt := checkers.UNKNOWN
	msg := fmt.Sprintf("%s: service does not exist.", opts.ServiceName)
	for _, s := range ss {
		if opts.ExcludeService != "" && strings.Contains(s.Name, opts.ExcludeService) {
			continue
		}
		if opts.Exact {
			if s.Name != opts.ServiceName {
				continue
			}
		} else {
			if !strings.Contains(s.Name, opts.ServiceName) {
				continue
			}
		}
		if s.State == "Running" {
			checkSt = checkers.OK
		} else {
			checkSt = checkers.CRITICAL
		}
		msg = fmt.Sprintf("%s: %s - %s", s.Name, s.Caption, s.State)
		break
	}

	return checkers.NewChecker(checkSt, msg)
}
