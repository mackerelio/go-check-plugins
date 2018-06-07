package mock

import (
	"github.com/mackerelio/checkers"
)


// Executer is dummy executer for testing
type Executer struct {
	CommandName string
	CommandArgs []string
	CommandResult string
	Status checkers.Status
	ParseResult string
}

// MakeCommandName is dummy command name maker for testing
func (e *Executer) MakeCommandName() string {
	return e.CommandName
}

// MakeCommandArgs is dummy command arguments maker for testing
func (e *Executer) MakeCommandArgs() []string {
	return e.CommandArgs
}

// Parse is dummy parse method for testing
func (e *Executer) Parse(result string) (checkers.Status, string) {
	e.CommandResult = result
	return e.Status, e.ParseResult
}
