package checkhttp

import (
	"fmt"
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	ckr := run([]string{"-u", "hoge"})
	assert.Equal(t, ckr.Status, checkers.CRITICAL, "chr.Status should be CRITICAL")
	assert.Equal(t, ckr.Message, `Get hoge: unsupported protocol scheme ""`, "something went wrong")
}

func TestNoCheckCertificate(t *testing.T) {
	ckr := run([]string{"--no-check-certificate", "-u", "hoge"})
	assert.Equal(t, ckr.Status, checkers.CRITICAL, "chr.Status should be CRITICAL")
	assert.Equal(t, ckr.Message, `Get hoge: unsupported protocol scheme ""`, "something went wrong")
}

func TestStatusRange(t *testing.T) {
	tests := []struct {
		args []string
		want checkers.Status
		err  bool
	}{
		{
			args: []string{"-s", "404=ok", "-u", "hoge"},
			want: checkers.CRITICAL,
		},
		{
			args: []string{"-s", "404=ok", "-u", "https://mackerel.io/404"},
			want: checkers.OK,
		},
		{
			args: []string{"-s", "401=ok", "-u", "https://mackerel.io/404"},
			want: checkers.WARNING,
		},
		{
			args: []string{"-s", "300-404=ok", "-u", "https://mackerel.io/404"},
			want: checkers.OK,
		},
		{
			args: []string{"-s", "200-300-404=ok", "-u", "http://example.com"},
			want: checkers.UNKNOWN,
		},
		{
			args: []string{"-s", "300=ok", "-s", "404=ok", "-u", "https://mackerel.io/404"},
			want: checkers.OK,
		},
		{
			args: []string{"-s", "300", "-s", "404=ok", "-u", "https://mackerel.io/404"},
			want: checkers.UNKNOWN,
		},
		{
			args: []string{"-s", "=ok", "-s", "404=ok", "-u", "https://mackerel.io/404"},
			want: checkers.UNKNOWN,
		},
		{
			args: []string{"-s", "=ok", "-s", "404-200=ok", "-u", "https://mackerel.io/404"},
			want: checkers.UNKNOWN,
		},
	}
	for _, tt := range tests {
		ckr := run(tt.args)
		assert.Equal(t, ckr.Status, tt.want, fmt.Sprintf("chr.Status wrong: %v", ckr.Status))
	}
}

func TestSourceIP(t *testing.T) {
	ckr := run([]string{"-u", "hoge", "-i", "1.2.3"})
	assert.Equal(t, ckr.Status, checkers.UNKNOWN, "chr.Status should be UNKNOWN")
}
