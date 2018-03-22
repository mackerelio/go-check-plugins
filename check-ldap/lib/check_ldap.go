package checkldap

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"gopkg.in/ldap.v2"
)

type checkLDAPOpts struct {
	Warning   float64 `short:"w" long:"warning" description:"Response time to result in warning status (seconds)" required:"true"`
	Critical  float64 `short:"c" long:"critical" descriptin:"Response time to result in critical status (seconds)" required:"true"`
	Host      string  `short:"H" long:"host" default:"localhost" description:"Hostname"`
	Port      string  `short:"p" long:"port" default:"389" description:"Port number"`
	Base      string  `short:"b" long:"base" description:"LDAP base" required:"true"`
	Attribute string  `short:"a" long:"attr" default:"(objectclass=*)" description:"LDAP attribute to search"`
	BindDN    string  `short:"D" long:"bind" description:"LDAP bind DN"`
	Password  string  `short:"P" long:"password" description:"LDAP password"`
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "LDAP"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	opts := checkLDAPOpts{}
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	lconn, err := ldap.Dial("tcp", net.JoinHostPort(opts.Host, opts.Port))
	if err != nil {
		return checkers.Unknown(err.Error())
	}
	defer lconn.Close()

	if err := lconn.Bind(opts.BindDN, opts.Password); err != nil {
		return checkers.Unknown(err.Error())
	}

	req := ldap.NewSearchRequest(
		opts.Base,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		opts.Attribute,
		nil,
		nil,
	)

	stTime := time.Now()

	_, err = lconn.Search(req)
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	elapsed := time.Since(stTime)

	msg := fmt.Sprintf("%.3f seconds response time", elapsed.Seconds())

	if elapsed.Seconds() > opts.Critical {
		return checkers.Critical(msg)
	} else if elapsed.Seconds() > opts.Warning {
		return checkers.Warning(msg)
	} else {
		return checkers.Ok(msg)
	}
}
