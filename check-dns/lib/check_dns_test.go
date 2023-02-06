package checkdns

import (
	"net"
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

func TestDnsStatus(t *testing.T) {
	ckr := Run([]string{"-H", "example.com", "-s", "8.8.8.8"})
	assert.Equal(t, checkers.OK, ckr.Status, "should be OK")
	t.Logf(ckr.Message)
	if !strings.Contains(ckr.Message, "status: NOERROR") {
		t.Errorf("status is not NOERROR")
	}
}
