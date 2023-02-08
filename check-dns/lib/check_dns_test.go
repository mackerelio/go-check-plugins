package checkdns

import (
	"net"
	"os"
	"strings"
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestNameServer(t *testing.T) {
	nameserver, err := adapterAddress()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf(nameserver)
	address := net.ParseIP(nameserver)
	if address == nil {
		t.Errorf("nameserver is invalid IP")
	}
}

func TestCheckDns(t *testing.T) {
	tests := []struct {
		args        []string
		want_status checkers.Status
		want_msg    string
	}{
		{
			[]string{"-H", "example.com"},
			checkers.OK,
			"status: NOERROR",
		},
		{
			[]string{"-H", "example.com", "--norec"},
			checkers.OK,
			"status: NOERROR",
		},
		{
			[]string{"-H", "exampleeeee.com"},
			checkers.CRITICAL,
			"status: NXDOMAIN",
		},
		{
			[]string{"-H", "example.com", "-s", "8.8.8.8"},
			checkers.OK,
			"status: NOERROR",
		},
		{
			[]string{"-H", "exampleeeee.com", "-s", "8.8.8.8"},
			checkers.CRITICAL,
			"status: NXDOMAIN",
		},
		{
			[]string{"-H", "exampleeeee.com", "-s", "8.8.8"},
			checkers.CRITICAL,
			"timeout",
		},
		{
			[]string{"-H", "jprs.co.jp", "-s", "202.11.16.49", "--norec"},
			checkers.OK,
			"status: NOERROR",
		},
		{
			[]string{"-H", "www.google.com", "-s", "202.11.16.49", "--norec"},
			checkers.CRITICAL,
			"status: REFUSED",
		},
	}

	for i, tt := range tests {
		t.Logf("=== Start #%d", i)
		opts, err := parseArgs(tt.args)
		if err != nil {
			t.Fatal(err)
		}
		// when runs without setting server in CI, status will be REFUSED
		if opts.Server == "" && os.Getenv("RUN_TEST_ON_GITHUB_ACTIONS") == "1" {
			continue
		}
		ckr := opts.run()

		assert.Equal(t, tt.want_status, ckr.Status)

		if !strings.Contains(ckr.Message, tt.want_msg) {
			t.Errorf("%s is not incleded in message", tt.want_msg)
		}
	}
}
