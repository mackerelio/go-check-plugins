package checkmasterha

import (
	"strconv"
	"strings"

	"github.com/mackerelio/checkers"
)

type replChecker struct {
	subcommand
	SecondsBehindMaster int `long:"seconds_behind_master" description:"seconds_behind_master option for masterha_check_repl"`
}

func (c replChecker) Execute(args []string) error {
	c.Executer = &c

	checker := c.executeAll()
	checker.Name = "MasterHA"
	checker.Exit()
	return nil
}

func (c replChecker) MakeCommandName() string {
	return "masterha_check_repl"
}

func (c replChecker) MakeCommandArgs() []string {
	args := make([]string, 0, c.ArgsLength())
	if c.SecondsBehindMaster > 0 {
		secondsBehindMaster := strconv.Itoa(c.SecondsBehindMaster)
		args = append(args, "--seconds_behind_master", secondsBehindMaster)
	}
	return args
}

func (c replChecker) ArgsLength() int {
	if c.SecondsBehindMaster > 0 {
		return 4
	}
	return 2
}

func (c replChecker) Parse(out string) (checkers.Status, string) {
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
