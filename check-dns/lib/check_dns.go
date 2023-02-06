package checkdns

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/miekg/dns"
)

type dnsOpts struct {
	Host   string `short:"H" long:"host" required:"true" description:"The name or address you want to query"`
	Server string `short:"s" long:"server" description:"DNS server you want to use for the lookup"`
	Port   int    `short:"p" long:"port" default:"53" description:"Port number you want to use"`
	Norec  bool   `long:"norec" description:"Set not recursive mode"`
}

// Do the plugin
func Do() {
	opts, err := parseArgs(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
	ckr := opts.run()
	ckr.Name = "DNS"
	ckr.Exit()
}

func parseArgs(args []string) (*dnsOpts, error) {
	opts := &dnsOpts{}
	_, err := flags.ParseArgs(opts, args)
	return opts, err
}

func (opts *dnsOpts) run() *checkers.Checker {
	var nameserver string
	var err error
	if opts.Server != "" {
		nameserver = opts.Server
	} else {
		nameserver, err = adapterAddress()
		if err != nil {
			return checkers.Critical(err.Error())
		}
	}
	nameserver = net.JoinHostPort(nameserver, strconv.Itoa(opts.Port))

	c := new(dns.Client)
	m := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			RecursionDesired: !opts.Norec,
			Opcode:           dns.OpcodeQuery,
		},
		Question: []dns.Question{{Name: dns.Fqdn(opts.Host), Qtype: dns.TypeA, Qclass: uint16(dns.ClassINET)}},
	}
	m.Id = dns.Id()

	r, _, err := c.Exchange(m, nameserver)
	if err != nil {
		return checkers.Critical(err.Error())
	}

	checkSt := checkers.OK
	if r.MsgHdr.Rcode != dns.RcodeSuccess {
		checkSt = checkers.CRITICAL
	}
	msg := fmt.Sprintf("%s", r)

	return checkers.NewChecker(checkSt, msg)
}
