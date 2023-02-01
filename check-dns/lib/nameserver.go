//go:build !windows

package checkdns

import (
	"net"
	"fmt"
	"github.com/miekg/dns"
)

func adapterAddress() (string, error) {
	conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return "", err
	}
	nameserver := conf.Servers[0]
	// ref: https://github.com/miekg/exdns/blob/d851fa434ad51cb84500b3e18b8aa7d3bead2c51/q/q.go#L148-L153
	// if the nameserver is from /etc/resolv.conf the [ and ] are already
	// added, thereby breaking net.ParseIP. Check for this and don't
	// fully qualify such a name
	if nameserver[0] == '[' && nameserver[len(nameserver)-1] == ']' {
		nameserver = nameserver[1 : len(nameserver)-1]
	}
	// ref: https://github.com/miekg/exdns/blob/d851fa434ad51cb84500b3e18b8aa7d3bead2c51/q/q.go#L154-L158
	if net.ParseIP(nameserver) == nil {
		nameserver = dns.Fqdn(nameserver)
	}
	if net.ParseIP(nameserver) == nil {
		return "", fmt.Errorf("invalid nameserver: %s", nameserver)
	}
	return nameserver, nil
}
