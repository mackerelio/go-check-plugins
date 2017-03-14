package checksolr

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type solrOpts struct {
	Host string `short:"H" long:"host" default:"localhost" description:"Hostname"`
	Port string `short:"p" long:"port" default:"8983" description:"Port"`
	Core string `short:"c" long:"core" required:"true" description:"Core"`
}

func (s solrOpts) createBaseURL() string {
	return fmt.Sprintf("http://%s:%s/solr/%s", s.Host, s.Port, s.Core)
}

var commands = map[string](func(solrOpts) *checkers.Checker){
	"ping": checkPing,
}

func separateSub(argv []string) (string, []string) {
	if len(argv) == 0 || strings.HasPrefix(argv[0], "-") {
		return "", argv
	}
	return argv[0], argv[1:]
}

// Do the plugin
func Do() {
	subCmd, argv := separateSub(os.Args[1:])
	fn, ok := commands[subCmd]
	if !ok {
		fmt.Println(`Usage:
  check-solr [subcommand] [OPTIONS]

SubCommands:`)
		for k := range commands {
			fmt.Printf("  %s\n", k)
		}
		os.Exit(1)
	}

	opts := solrOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	psr.Usage = fmt.Sprintf("%s [OPTIONS]", subCmd)
	_, err := psr.ParseArgs(argv)
	if err != nil {
		os.Exit(1)
	}

	ckr := fn(opts)
	ckr.Name = fmt.Sprintf("Solr %s", strings.Title(subCmd))
	ckr.Exit()
}
