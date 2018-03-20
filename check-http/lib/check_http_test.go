package checkhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	ckr := Run([]string{"-u", "hoge"})
	assert.Equal(t, ckr.Status, checkers.CRITICAL, "chr.Status should be CRITICAL")
	assert.Equal(t, ckr.Message, `Get hoge: unsupported protocol scheme ""`, "something went wrong")
}

func TestNoCheckCertificate(t *testing.T) {
	ckr := Run([]string{"--no-check-certificate", "-u", "hoge"})
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
		ckr := Run(tt.args)
		assert.Equal(t, ckr.Status, tt.want, fmt.Sprintf("chr.Status wrong: %v", ckr.Status))
	}
}

func TestSourceIP(t *testing.T) {
	ckr := Run([]string{"-u", "hoge", "-i", "1.2.3"})
	assert.Equal(t, ckr.Status, checkers.UNKNOWN, "chr.Status should be UNKNOWN")
}

func TestExpectedContent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	testCases := []struct {
		regexp string
		status checkers.Status
	}{
		{
			regexp: "Hello, client",
			status: checkers.OK,
		},
		{
			regexp: "Wrong response",
			status: checkers.CRITICAL,
		},
		{
			regexp: "Hel.*",
			status: checkers.OK,
		},
		{
			regexp: "clientt?",
			status: checkers.OK,
		},
		{
			regexp: "???",
			status: checkers.UNKNOWN,
		},
	}

	for i, tc := range testCases {
		ckr := Run([]string{"-u", ts.URL, "-p", tc.regexp})
		assert.Equal(t, ckr.Status, tc.status, "#%d: Status should be %s", i, tc.status)
	}

	//ckr := Run([]string{"-u", ts.URL, "-p", "Hello, client"})
	//assert.Equal(t, ckr.Status, checkers.OK, "chr.Status should be OK")

	//ckr = Run([]string{"-u", ts.URL, "-p", "Wrong response"})
	//assert.Equal(t, ckr.Status, checkers.CRITICAL, "chr.Status should be CRITICAL")
}
