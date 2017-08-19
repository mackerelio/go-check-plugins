package checkmasterha

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

// Do the plugin
func Do() {
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

type executer interface {
	MakeCommandName() string
	MakeCommandArgs() []string
	Parse(string) (checkers.Status, string)
}

type subcommand struct {
	Config    string `short:"c" long:"conf" description:"target config file"`
	ConfigDir string `long:"confdir" default:"/usr/local/masterha/conf" description:"config directory"`
	All       bool   `short:"a" long:"all" description:"use all config file for target"`
	Executer  executer
}

func (c subcommand) ConfigFiles() ([]string, error) {
	if c.All {
		files, err := ioutil.ReadDir("/service")
		if err != nil {
			return nil, err
		}

		configFiles := make([]string, 0, len(files))
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "masterha_") {
				configFile := c.ConfigDir + "/" + file.Name()[9:] + ".cnf"
				configFiles = append(configFiles, configFile)
			}
		}
		return configFiles, nil
	}

	configFiles := []string{c.Config}
	return configFiles, nil
}

func (c subcommand) MakeCommandName() string {
	return c.Executer.MakeCommandName()
}

func (c subcommand) MakeCommandArgs() []string {
	args := c.Executer.MakeCommandArgs()
	args = append(args, "--conf", c.Config)
	return args
}

func (c subcommand) Parse(result string) (checkers.Status, string) {
	return c.Executer.Parse(result)
}

func (c subcommand) executeAll() *checkers.Checker {
	var err error
	checker := checkers.NewChecker(checkers.UNKNOWN, "No target")

	configFiles, err := c.ConfigFiles()
	if err != nil {
		checker.Status  = checkers.UNKNOWN
		checker.Message = err.Error()
		return checker
	}

	for _, config := range configFiles {
		checker = c.execute(config)
		if checker.Status != checkers.OK {
			return checker
		}
	}

	return checker
}

func (c subcommand) execute(config string) *checkers.Checker {
	c.Config = config
	name := c.MakeCommandName()
	args := c.MakeCommandArgs()

	cmd := exec.Command(name, args...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	var failure bool
	err := cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		failure = true
	} else if err != nil {
		checker := checkers.NewChecker(checkers.UNKNOWN, err.Error())
		return checker
	}

	result, msg := c.Parse(buf.String())
	if failure && result == checkers.UNKNOWN {
		result = checkers.WARNING
	}
	checker := checkers.NewChecker(result, msg)
	return checker
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
