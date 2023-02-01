package checkdns

import (
	"github.com/miekg/dns"
)

func adapterAddress() (string, error) {
	conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return "", err
	}
	return conf.Servers[0], nil
}
