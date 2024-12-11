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
		t.Error(err)
	}
	t.Log(nameserver)
	address := net.ParseIP(nameserver)
	if address == nil {
		t.Error("nameserver is invalid IP")
	}
}

func TestCheckDns(t *testing.T) {
	tests := []struct {
		args        []string
		want_status checkers.Status
		want_msg    []string
	}{
		{
			[]string{"-H", "a.root-servers.net"},
			checkers.OK,
			[]string{"status: NOERROR"},
		},
		{
			[]string{"-H", "a.root-servers.net", "--norec"},
			checkers.OK,
			[]string{"status: NOERROR"},
		},
		{
			[]string{"-H", "a.root-servers.invalid"},
			checkers.CRITICAL,
			[]string{"status: NXDOMAIN"},
		},
		{
			[]string{"-H", "a.root-servers.net", "-s", "8.8.8.8"},
			checkers.OK,
			[]string{"status: NOERROR"},
		},
		{
			[]string{"-H", "a.root-servers.invalid", "-s", "8.8.8.8"},
			checkers.CRITICAL,
			[]string{"status: NXDOMAIN"},
		},
		{
			[]string{"-H", "a.root-servers.net", "-s", "8.8.8"},
			checkers.CRITICAL,
			[]string{""},
		},
		{
			[]string{"-H", "a.root-servers.net", "-s", "8.8.8.8", "-q", "AAAA"},
			checkers.OK,
			[]string{"status: NOERROR", "AAAA"},
		},
		{
			[]string{"-H", "a.root-servers.net", "-s", "8.8.8.8", "-q", "AAA"},
			checkers.CRITICAL,
			[]string{"AAA is invalid query type"},
		},
		{
			[]string{"-H", "a.root-servers.net", "-s", "8.8.8.8", "-e", "198.41.0.4"},
			checkers.OK,
			[]string{"status: NOERROR", "198.41.0.4"},
		},
		{
			[]string{"-H", "a.root-servers.net", "-s", "8.8.8.8", "-q", "AAAA", "--expected-string", "2001:503:ba3e::2:30"},
			checkers.OK,
			[]string{"status: NOERROR", "2001:503:ba3e::2:30"},
		},
		{
			[]string{"-H", "a.root-servers.net", "-s", "8.8.8.8", "-q", "TXT", "-e", ""},
			checkers.CRITICAL,
			[]string{"is not supported query type. Only A, AAAA is supported for expectation."},
		},
		{
			[]string{"-H", "a.root-servers.net", "-s", "8.8.8.8", "-e", "198.41.0.3"},
			checkers.CRITICAL,
			[]string{"status: NOERROR", "198.41.0.4"},
		},
		{
			[]string{"-H", "a.root-servers.invalid", "-s", "8.8.8.8", "-e", "198.41.0.4"},
			checkers.CRITICAL,
			[]string{"status: NXDOMAIN"},
		},
		{
			[]string{"-H", "a.root-servers.net", "-s", "8.8.8.8", "-e", "198.41.0.3", "-e", "198.41.0.4"},
			checkers.WARNING,
			[]string{"status: NOERROR", "198.41.0.4"},
		},
		{
			[]string{"-H", "a.root-servers.net", "-s", "8.8.8.8", "-e", "198.41.0.3", "-e", "   198.41.0.4   "},
			checkers.CRITICAL,
			[]string{"status: NOERROR", "198.41.0.4"},
		},
	}

	for i, tt := range tests {
		t.Logf("=== Start #%d", i)
		opts, err := parseArgs(tt.args)
		if err != nil {
			t.Fatal(err)
		}
		// when runs without setting server in CI, status will be REFUSED
		if opts.Server == "" && os.Getenv("CI") == "true" {
			continue
		}
		ckr := opts.run()

		assert.Equal(t, tt.want_status, ckr.Status)

		for _, want := range tt.want_msg {
			if !strings.Contains(ckr.Message, want) {
				t.Errorf("%s is not incleded in message", want)
			}
		}
	}
}
