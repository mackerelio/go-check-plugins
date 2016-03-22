package main

import (
	"strings"

	"github.com/mackerelio/checkers"
)

type replChecker struct {
	Config string `short:"c" long:"conf" required:"true" description:"target config file"`
}

func (c replChecker) Execute(args []string) error {
	checker, err := executeSubcommand(c)
	if err != nil {
		return err
	}
	checker.Name = "MasterHA"
	checker.Exit()
	return nil
}

func (c replChecker) makeCommandName() string {
	return "masterha_check_repl"
}

func (c replChecker) makeCommandArgs() []string {
	args := [2]string{"--conf", c.Config}
	return args[:]
}

func (c replChecker) parse(out string) (checkers.Status, string) {
	lines := extractNonEmptyLines(strings.Split(out, "\n"))
	lastLine := lines[len(lines)-1]
	if strings.Contains(lastLine, "MySQL Replication Health is OK.") {
		return checkers.OK, lastLine
	} else if strings.Contains(lastLine, "MySQL Replication Health is NOT OK!") {
		return checkers.CRITICAL, lastLine
	} else {
		msg := extractErrorMsg(out)
		return checkers.UNKNOWN, msg
	}
}
