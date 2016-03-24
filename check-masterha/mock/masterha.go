package mock

import (
	"github.com/mackerelio/checkers"
)

type Executer struct {
	CommandName string
	CommandArgs []string
	CommandResult string
	Status checkers.Status
	ParseResult string
}

func (e *Executer) MakeCommandName() string {
	return e.CommandName
}

func (e *Executer) MakeCommandArgs() []string {
	return e.CommandArgs
}

func (e *Executer) Parse(result string) (checkers.Status, string) {
	e.CommandResult = result
	return e.Status, e.ParseResult
}
