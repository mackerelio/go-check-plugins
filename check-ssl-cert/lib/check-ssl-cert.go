package checksslcert

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type certOpts struct {
	Host     string `short:"H" long:"host" required:"true" description:"Host name"`
	Port     int    `short:"p" long:"port" default:"443" description:"Port number"`
	Warning  int    `short:"w" long:"warning" value-name:"days" default:"30" description:"The warning threshold in days before expiry"`
	Critical int    `short:"c" long:"critical" value-name:"days" default:"14" description:"The critical threshold in days before expiry"`
}

func parseArgs(args []string) (*certOpts, error) {
	opts := &certOpts{}
	_, err := flags.ParseArgs(opts, args)
	return opts, err
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "SSL"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	opts, err := parseArgs(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	cert, err := getCert(addr)
	if err != nil {
		return checkers.Critical(err.Error())
	}
	expiry := cert.NotAfter
	dur := expiry.Sub(time.Now())

	chkSt := checkers.OK
	days := int(dur.Hours() / 24)
	dayMsg := fmt.Sprintf("%d ", days)
	if days == 1 {
		dayMsg += "day"
	} else {
		dayMsg += "days"
	}
	msg := fmt.Sprintf("Certificate '%s' expires in %s (%s)", addr, dayMsg, expiry)
	if dur < time.Duration(opts.Warning)*time.Hour*24 {
		chkSt = checkers.WARNING
	}
	if dur < time.Duration(opts.Critical)*time.Hour*24 {
		chkSt = checkers.CRITICAL
	}
	return checkers.NewChecker(chkSt, msg)
}

func getCert(addr string) (*x509.Certificate, error) {
	conn, err := tls.Dial("tcp", addr, &tls.Config{})
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates

	if len(certs) < 1 {
		return nil, fmt.Errorf("no certifiations are available")
	}
	return certs[0], err
}
