package main

import (
	"strings"

	"github.com/mackerelio/checkers"
)

type sshChecker struct {
	subcommand
}

func (c sshChecker) Execute(args []string) error {
	c.Executer = &c

	checker, err := c.executeAll()
	if err != nil {
		return err
	}

	checker.Name = "MasterHA"
	checker.Exit()
	return nil
}

func (c sshChecker) makeCommandName() string {
	return "masterha_check_ssh"
}

func (c sshChecker) makeCommandArgs() []string {
	return make([]string, 0, 2)
}

func (c sshChecker) parse(out string) (checkers.Status, string) {
	lines := extractNonEmptyLines(strings.Split(out, "\n"))
	lastLine := lines[len(lines)-1]
	if strings.Contains(lastLine, "All SSH connection tests passe") {
		return checkers.OK, lastLine
	} else if strings.Contains(lastLine, "SSH Configuration Check Failed!") {
		return checkers.OK, lastLine
	} else {
		msg := extractErrorMsg(out)
		return checkers.UNKNOWN, msg
	}
}
