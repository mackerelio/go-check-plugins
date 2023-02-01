package checkdns

import (
	"net"
	"testing"
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
