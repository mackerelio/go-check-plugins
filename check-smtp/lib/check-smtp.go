package checksmtp

import (
	"os"

	"net/smtp"

	"fmt"

	"crypto/tls"

	"net"

	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var opts struct {
	Host     string  `short:"H" long:"host" default:"localhost" description:"Hostname"`
	Port     string  `short:"p" long:"port" default:"25" description:"Port"`
	FQDN     string  `short:"F" long:"fqdn" description:"FQDN used for HELO"`
	SMTPS    bool    `short:"s" long:"smtps" description:"Use SMTP over TLS"`
	StartTLS bool    `short:"S" long:"starttls" description:"Use STARTTLS"`
	Auth     string  `short:"A" long:"authmech" description:"SMTP AUTH Authentication Mechanisms (only PLAIN supported)"`
	User     string  `short:"U" long:"authuser" description:"SMTP AUTH username"`
	Password string  `short:"P" long:"authpassword" description:"SMTP AUTH password"`
	Warning  float64 `short:"w" long:"warning" description:"Warning threshold (sec)"`
	Critical float64 `short:"c" long:"critical" description:"Critical threshold (sec)"`
	Timeout  int     `short:"t" long:"timeout" default:"10" description:"Timeout (sec)"`
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "SMTP"
	ckr.Exit()
}

func makeConn(host, port string, timeout int, isSMTPS bool, tlsConfig *tls.Config) (net.Conn, error) {
	d := net.Dialer{Timeout: time.Duration(timeout) * time.Second}
	if isSMTPS {
		return tls.DialWithDialer(&d, "tcp", fmt.Sprintf("%s:%s", host, port), tlsConfig)
	}
	return d.Dial("tcp", fmt.Sprintf("%s:%s", opts.Host, opts.Port))
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	if opts.Warning == 0 && opts.Critical == 0 {
		return checkers.Unknown("require threshold option (warning or critical)")
	}

	if opts.Auth != "" && opts.Auth != "PLAIN" {
		return checkers.Unknown("invalid SMTP AUTH Authentication Mechanisms (only PLAIN supported)")
	}

	fqdn := opts.FQDN
	if fqdn == "" {
		localHostName, err := os.Hostname()
		if err != nil {
			return checkers.Unknown(err.Error())
		}
		fqdn = localHostName
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         opts.Host,
	}

	conn, err := makeConn(opts.Host, opts.Port, opts.Timeout, opts.SMTPS, tlsConfig)
	if err != nil {
		return checkers.Critical(err.Error())
	}

	stTime := time.Now()

	c, err := smtp.NewClient(conn, opts.Host)
	if err != nil {
		return checkers.Critical(err.Error())
	}
	defer c.Quit()

	if err := c.Hello(fqdn); err != nil {
		return checkers.Critical(err.Error())
	}

	if opts.StartTLS {
		if err := c.StartTLS(tlsConfig); err != nil {
			return checkers.Critical(err.Error())
		}
	}

	if opts.Auth == "PLAIN" {
		auth := smtp.PlainAuth("", opts.User, opts.Password, opts.Host)
		if err := c.Auth(auth); err != nil {
			return checkers.Critical(err.Error())
		}
	}

	if err := c.Mail(""); err != nil {
		c.Reset()
		return checkers.Critical(err.Error())
	}
	c.Reset()

	elapsed := time.Since(stTime)

	msg := fmt.Sprintf("%.3f seconds response time", elapsed.Seconds())

	if opts.Critical != 0 && elapsed.Seconds() > opts.Critical {
		return checkers.Critical(msg)
	} else if opts.Warning != 0 && elapsed.Seconds() > opts.Warning {
		return checkers.Warning(msg)
	} else {
		return checkers.Ok(msg)
	}
}
