package checkmasterha

import (
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestMakeStatusCommandName(t *testing.T) {
	name := statusSubcommand.MakeCommandName()
	assert.Equal(t, "masterha_check_status", name)
}

func TestMakeStatusCommandArgs(t *testing.T) {
	args := statusSubcommand.MakeCommandArgs()
	assert.Equal(t, 0, len(args))
}

func TestParseStatusSuccess(t *testing.T) {
	body := `db001 (pid:1111) is running(0:PING_OK), master:XX.XX.XX.XX
`
	status, msg := statusSubcommand.Parse(body)
	assert.Equal(t, checkers.OK, status)
	assert.Equal(t, "running(0:PING_OK)", msg)
}

func TestParseStatusFailure(t *testing.T) {
	body := `db001 is stopped(2:NOT_RUNNING).
`
	status, msg := statusSubcommand.Parse(body)
	assert.Equal(t, checkers.CRITICAL, status)
	assert.Equal(t, "db001 is stopped(2:NOT_RUNNING).", msg)
}

