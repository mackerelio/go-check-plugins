package checkping

import (
	"os"

	"github.com/mackerelio/checkers"
)

func run(args []string) *checkers.Checker {
	return checkers.NewChecker(checkers.OK, "")
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "Ping"
	ckr.Exit()
}
