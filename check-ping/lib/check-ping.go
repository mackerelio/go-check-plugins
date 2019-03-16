package checkping

import (
	"net"
	"os"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	ping "github.com/tatsushid/go-fastping"
)

var opts struct {
	Host     string `long:"host" short:"H" description:"check target IP Address"`
	Count    int    `long:"count" short:"n" default:"1" description:"sending (and receiving) count ping packets"`
	WaitTime int    `long:"wait-time" short:"w" default:"1000" description:"wait time, Max RTT(ms)"`
}

func run(args []string) *checkers.Checker {
	var parser = flags.NewParser(&opts, flags.Default)
	_, err := parser.ParseArgs(args)
	if err != nil {
		os.Exit(1)
	}
	if opts.Host == "" {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	p := ping.NewPinger()
	netProto := "ip4:icmp"
	if strings.Index(opts.Host, ":") != -1 {
		netProto = "ip6:ipv6-icmp"
	}

	ra, err := net.ResolveIPAddr(netProto, opts.Host)
	if err != nil {
		os.Exit(1)
	}
	p.AddIPAddr(ra)

	status := checkers.CRITICAL
	p.MaxRTT = time.Millisecond * time.Duration(opts.WaitTime)
	p.OnRecv = func(_ *net.IPAddr, _ time.Duration) {
		status = checkers.OK
	}

	for i := 0; i < opts.Count; i++ {
		err := p.Run()
		if err != nil {
			return checkers.NewChecker(status, err.Error())
		}
	}

	return checkers.NewChecker(status, "")
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "Ping"
	ckr.Exit()
}
