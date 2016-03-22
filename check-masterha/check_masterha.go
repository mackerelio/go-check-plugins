package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}

type options struct {
	Status statusChecker `command:"status" description:"check to masterha_check_status"`
	Repl   replChecker   `command:"repl"   description:"check to masterha_check_repl"`
	SSH    sshChecker    `command:"ssh"    description:"check to masterha_check_ssh"`
}

type subcommand interface {
	makeCommandName() string
	makeCommandArgs() []string
	parse(string) (checkers.Status, string)
}

func executeSubcommand(c subcommand) (*checkers.Checker, error) {
	name := c.makeCommandName()
	args := c.makeCommandArgs()

	cmd := exec.Command(name, args...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	var failure bool
	err := cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		failure = true
	} else if err != nil {
		return nil, err
	}

	result, msg := c.parse(buf.String())
	if failure && result == checkers.UNKNOWN {
		result = checkers.WARNING
	}
	checker := checkers.NewChecker(result, msg)
	return checker, nil
}

func extractNonEmptyLines(lines []string) []string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

func extractErrorMsg(msg string) string {
	var errors []string
	for _, line := range strings.Split(msg, "\n") {
		if strings.Contains(line, "[error]") {
			errors = append(errors, line)
		}
	}

	if len(errors) == 0 {
		return msg
	}

	return strings.Join(errors, "\n")
}
