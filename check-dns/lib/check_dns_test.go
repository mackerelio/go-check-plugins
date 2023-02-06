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
		local_only  bool
	}{
		{
			[]string{"-H", "example.com"},
			checkers.OK,
			"status: NOERROR",
			true,
		},
		{
			[]string{"-H", "example.com", "--norec"},
			checkers.OK,
			"status: NOERROR",
			true,
		},
		{
			[]string{"-H", "exampleeeee.com"},
			checkers.CRITICAL,
			"status: NXDOMAIN",
			true,
		},
		{
			[]string{"-H", "example.com", "-s", "8.8.8.8"},
			checkers.OK,
			"status: NOERROR",
			false,
		},
		{
			[]string{"-H", "exampleeeee.com", "-s", "8.8.8.8"},
			checkers.CRITICAL,
			"status: NXDOMAIN",
			false,
		},
	}

	for i, tt := range tests {
		t.Logf("=== Start #%d", i)
		// when runs without setting server in CI, status will be REFUSED
		if tt.local_only && os.Getenv("RUN_TEST_ON_GITHUB_ACTIONS") == "1" {
			continue
		}
		opts, err := parseArgs(tt.args)
		if err != nil {
			t.Fatal(err)
		}
		ckr := opts.run()

		assert.Equal(t, tt.want_status, ckr.Status)

		if tt.want_msg != "" {
			if !strings.Contains(ckr.Message, tt.want_msg) {
				t.Errorf("%s is not incleded in message", tt.want_msg)
			}
		}
	}
}
