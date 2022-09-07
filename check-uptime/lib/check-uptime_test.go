package checkuptime

import (
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestCheckUptime(t *testing.T) {
	tests := []struct {
		args []string
		want checkers.Status
	}{
		{
			[]string{},
			checkers.OK,
		},
		{
			[]string{"-w", "-1"},
			checkers.OK,
		},
		{
			[]string{"-W", "-1"},
			checkers.WARNING,
		},
		{
			[]string{"-c", "-1"},
			checkers.OK,
		},
		{
			[]string{"-C", "-1"},
			checkers.CRITICAL,
		},
		{
			[]string{"-W", "-1", "-C", "-1"},
			checkers.CRITICAL,
		},
		{
			[]string{"--warn-under", "-1"},
			checkers.OK,
		},
		{
			[]string{"--warn-over", "-1"},
			checkers.WARNING,
		},
	}
	for _, tt := range tests {
		opts, err := parseArgs(tt.args)
		if err != nil {
			t.Fatal(err)
		}
		ckr := opts.run()
		assert.Equal(t, tt.want, ckr.Status)
	}
}
