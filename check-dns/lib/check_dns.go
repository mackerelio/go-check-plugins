package checkdns

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/miekg/dns"
)

type dnsOpts struct {
	Host           string   `short:"H" long:"host" required:"true" description:"The name or address you want to query"`
	Server         string   `short:"s" long:"server" description:"DNS server you want to use for the lookup"`
	Port           int      `short:"p" long:"port" default:"53" description:"Port number you want to use"`
	QueryType      string   `short:"q" long:"querytype" default:"A" description:"DNS record query type where TYPE =(A, AAAA, TXT, MX, CNAME)"`
	QueryClass     string   `short:"c" long:"queryclass" default:"IN" description:"DNS record class type where TYPE =(IN, CS, CH, HS, NONE, ANY)"`
	Norec          bool     `long:"norec" description:"Set not recursive mode"`
	ExpectedString []string `short:"e" long:"expected-string" description:"The string you expect the DNS server to return. If multiple responses are returned at once, you have to specify whole string"`
}

var supportedQueryType = map[string]int{"A": 1, "AAAA": 1, "TXT": 1, "MX": 1, "CNAME": 1}

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

	_, ok := supportedQueryType[strings.ToUpper(opts.QueryType)]
	if !ok {
		return checkers.Critical(fmt.Sprintf("%s is not supported query type", opts.QueryType))
	}
	queryType := dns.StringToType[strings.ToUpper(opts.QueryType)]
	queryClass, ok := dns.StringToClass[strings.ToUpper(opts.QueryClass)]
	if !ok {
		return checkers.Critical(fmt.Sprintf("%s is invalid query class", opts.QueryClass))
	}

	c := new(dns.Client)
	m := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			RecursionDesired: !opts.Norec,
			Opcode:           dns.OpcodeQuery,
		},
		Question: []dns.Question{{Name: dns.Fqdn(opts.Host), Qtype: queryType, Qclass: uint16(queryClass)}},
	}
	m.Id = dns.Id()

	r, _, err := c.Exchange(m, nameserver)
	if err != nil {
		return checkers.Critical(err.Error())
	}

	checkSt := checkers.OK
	/**
	  if DNS server return 1.1.1.1, 2.2.2.2
		1: -e 1.1.1.1 -e 2.2.2.2            -> OK
		2: -e 1.1.1.1 -e 2.2.2.2 -e 3.3.3.3 -> WARNING
		3: -e 1.1.1.1 -e                    -> WARNING
		4: -e 1.1.1.1 -e 3.3.3.3            -> WARNING
		5: -e 3.3.3.3 -e                    -> CRITICAL
		6: -e 3.3.3.3 -e 4.4.4.4 -e 5.5.5.5 -> CRITICAL
	**/
	if len(opts.ExpectedString) != 0 {
		match := 0
		for _, expectedString := range opts.ExpectedString {
			for _, answer := range r.Answer {
				var anserWithoutHeader string
				expectMatch := expectedString
				switch t := answer.(type) {
				case *dns.A:
					anserWithoutHeader = t.A.String()
				case *dns.AAAA:
					anserWithoutHeader = t.AAAA.String()
				case *dns.TXT:
					anserWithoutHeader = sprintTxt(t.Txt)
					// " is added by sprintTxt
					expectMatch = "\"" + expectedString + "\""
				case *dns.MX:
					anserWithoutHeader = strconv.Itoa(int(t.Preference)) + " " + sprintName(t.Mx)
				case *dns.CNAME:
					anserWithoutHeader = sprintName(t.Target)
				default:
					return checkers.Critical(fmt.Sprintf("%s is not supported query type", opts.QueryType))
				}
				if anserWithoutHeader == expectMatch {
					match += 1
				}
			}
		}
		if match == len(r.Answer) {
			if len(opts.ExpectedString) == len(r.Answer) { // case 1
				checkSt = checkers.OK
			} else { // case 2
				checkSt = checkers.WARNING
			}
		} else {
			if match > 0 { // case 3,4
				checkSt = checkers.WARNING
			} else { // case 5,6
				checkSt = checkers.CRITICAL
			}
		}
	}

	if r.MsgHdr.Rcode != dns.RcodeSuccess {
		checkSt = checkers.CRITICAL
	}

	msg := fmt.Sprintf("HEADER-> %s\n", r.MsgHdr.String())
	for _, answer := range r.Answer {
		msg += fmt.Sprintf("ANSWER-> %s\n", answer)
	}

	return checkers.NewChecker(checkSt, msg)
}
