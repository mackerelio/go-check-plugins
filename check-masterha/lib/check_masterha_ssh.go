package checkmasterha

import (
	"strings"

	"github.com/mackerelio/checkers"
)

type sshChecker struct {
	subcommand
}

func (c sshChecker) Execute(args []string) error {
	c.Executer = &c

	checker := c.executeAll()
	checker.Name = "MasterHA"
	checker.Exit()
	return nil
}

func (c sshChecker) MakeCommandName() string {
	return "masterha_check_ssh"
}

func (c sshChecker) MakeCommandArgs() []string {
	return make([]string, 0, 2)
}

func (c sshChecker) Parse(out string) (checkers.Status, string) {
	lines := extractNonEmptyLines(strings.Split(out, "\n"))
	lastLine := lines[len(lines)-1]
	if strings.Contains(lastLine, "All SSH connection tests passe") {
		return checkers.OK, lastLine
	} else if strings.Contains(lastLine, "SSH Configuration Check Failed!") {
		return checkers.CRITICAL, lastLine
	} else {
		msg := extractErrorMsg(out)
		return checkers.UNKNOWN, msg
	}
}
