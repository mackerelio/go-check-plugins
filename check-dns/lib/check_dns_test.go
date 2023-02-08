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
		want_msg    []string
	}{
		{
			[]string{"-H", "example.com"},
			checkers.OK,
			[]string{"status: NOERROR"},
		},
		{
			[]string{"-H", "example.com", "--norec"},
			checkers.OK,
			[]string{"status: NOERROR"},
		},
		{
			[]string{"-H", "exampleeeee.com"},
			checkers.CRITICAL,
			[]string{"status: NXDOMAIN"},
		},
		{
			[]string{"-H", "example.com", "-s", "8.8.8.8"},
			checkers.OK,
			[]string{"status: NOERROR"},
		},
		{
			[]string{"-H", "exampleeeee.com", "-s", "8.8.8.8"},
			checkers.CRITICAL,
			[]string{"status: NXDOMAIN"},
		},
		{
			[]string{"-H", "exampleeeee.com", "-s", "8.8.8"},
			checkers.CRITICAL,
			[]string{""},
		},
		{
			[]string{"-H", "jprs.co.jp", "-s", "202.11.16.49", "--norec"},
			checkers.OK,
			[]string{"status: NOERROR"},
		},
		{
			[]string{"-H", "www.google.com", "-s", "202.11.16.49", "--norec"},
			checkers.CRITICAL,
			[]string{"status: REFUSED"},
		},
		{
			[]string{"-H", "example.com", "-s", "8.8.8.8", "-q", "AAAA"},
			checkers.OK,
			[]string{"status: NOERROR", "AAAA"},
		},
		{
			[]string{"-H", "example.com", "-s", "8.8.8.8", "-q", "AAA"},
			checkers.CRITICAL,
			[]string{"AAA is invalid queryType"},
		},
		{
			[]string{"-H", "example.com", "-s", "8.8.8.8", "-c", "IN"},
			checkers.OK,
			[]string{"status: NOERROR"},
		},
		{
			[]string{"-H", "example.com", "-s", "8.8.8.8", "-c", "INN"},
			checkers.CRITICAL,
			[]string{"INN is invalid queryClass"},
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

		for _, want := range tt.want_msg {
			if !strings.Contains(ckr.Message, want) {
				t.Errorf("%s is not incleded in message", want)
			}
		}
	}
}
