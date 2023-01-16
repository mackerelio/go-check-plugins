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
	Host               string `short:"H" long:"host" required:"true" description:"Host name"`
	Port               int    `short:"p" long:"port" default:"443" description:"Port number"`
	Warning            int    `short:"w" long:"warning" value-name:"days" default:"30" description:"The warning threshold in days before expiry"`
	Critical           int    `short:"c" long:"critical" value-name:"days" default:"14" description:"The critical threshold in days before expiry"`
	CAfile             string `long:"ca-file" description:"A CA Cert file to use for server authentication"`
	CertFile           string `long:"cert-file" description:"A Cert file to use for client authentication"`
	KeyFile            string `long:"key-file" description:"A Key file to use for client authentication"`
	NoCheckCertificate bool   `long:"no-check-certificate" description:"Do not check certificate"`
}

func parseArgs(args []string) (*certOpts, error) {
	opts := &certOpts{}
	_, err := flags.ParseArgs(opts, args)
	return opts, err
}

// Do the plugin
func Do() {
	ckr := Run(os.Args[1:])
	ckr.Name = "SSL"
	ckr.Exit()
}

func Run(args []string) *checkers.Checker {
	opts, err := parseArgs(args)
	if err != nil {
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	cert, err := getCert(addr, opts)
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

func getCert(addr string, opts *certOpts) (*x509.Certificate, error) {
	tlsConfig := &tls.Config{}
	if opts.CAfile != "" {
		caCert, err := os.ReadFile(opts.CAfile)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
		tlsConfig.BuildNameToCertificate()
	}
	if opts.CertFile != "" && opts.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(opts.CertFile, opts.KeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	if opts.NoCheckCertificate {
		tlsConfig.InsecureSkipVerify = true
	}
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates

	if len(certs) < 1 {
		return nil, fmt.Errorf("no certifiations are available")
	}
	cert := certs[0]
	for _, c := range certs[1:] {
		if c.NotAfter.Before(cert.NotAfter) {
			cert = c
		}
	}
	return cert, err
}
