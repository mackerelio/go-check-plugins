package checkcloudwatchlogs

import (
	"os"

	"github.com/mackerelio/checkers"
)

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "CloudWatch Logs"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	return checkers.NewChecker(checkers.OK, "ok")
}
