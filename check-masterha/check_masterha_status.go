package main

import (
	"strings"

	"github.com/mackerelio/checkers"
)

type statusChecker struct {
	Config string `short:"c" long:"conf" default:"" description:"target config file"`
	All    bool   `short:"a" long:"all" description:"use all config file for target"`
}

func (c statusChecker) Execute(args []string) error {
	checker, err := executeSubcommand(c)
	if err != nil {
		return err
	}
	checker.Name = "MasterHA"
	checker.Exit()
	return nil
}

func (c statusChecker) makeCommandName() string {
	return "masterha_check_status"
}

func (c statusChecker) makeCommandArgs() []string {
	args := make([]string, 0, 2)
	if c.All {
		args = append(args, "--all")
	} else if c.Config != "" {
		args = append(args, "--conf", c.Config)
	}
	return args
}

func (c statusChecker) parse(out string) (checkers.Status, string) {
	lines := strings.Split(out, "\n")
	errors := make([]string, 0, 0)

	for _, line := range lines {
		if line != "" && !strings.Contains(line, "running(0:PING_OK)") {
			errors = append(errors, line)
		}
	}
	if len(errors) == 0 {
		return checkers.OK, "running(0:PING_OK)"
	}

	msg := strings.Join(errors, "\n")
	return checkers.CRITICAL, msg
}
