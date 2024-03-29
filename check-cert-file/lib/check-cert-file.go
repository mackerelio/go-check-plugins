package checkcertfile

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type certOpts struct {
	CertFile string `short:"f" long:"file" required:"true" description:"cert file name"`
	Crit     int64  `short:"c" long:"critical" default:"14" description:"The critical threshold in days before expiry"`
	Warn     int64  `short:"w" long:"warning" default:"30" description:"The threshold in days before expiry"`
}

// Do the plugin
func Do() {
	ckr := Run(os.Args[1:])
	ckr.Name = "CERT Expiry"
	ckr.Exit()
}

func parseArgs(args []string) (*certOpts, error) {
	opts := &certOpts{}
	_, err := flags.ParseArgs(opts, args)
	return opts, err
}

func Run(args []string) *checkers.Checker {
	opts, err := parseArgs(args)
	if err != nil {
		os.Exit(1)
	}

	cfByte, err := os.ReadFile(opts.CertFile)
	if err != nil {
		return checkers.Critical(err.Error())
	}
	cfBlock, _ := pem.Decode(cfByte)
	cfCrt, err := x509.ParseCertificate(cfBlock.Bytes)

	if err != nil {
		return checkers.Critical(err.Error())
	}

	cfDaysRemaining := int64(cfCrt.NotAfter.Sub(time.Now().UTC()).Hours() / 24)
	checkSt := checkers.OK
	msg := fmt.Sprintf("%d days remaining", cfDaysRemaining)

	if cfDaysRemaining < opts.Crit {
		checkSt = checkers.CRITICAL
	} else if cfDaysRemaining < opts.Warn {
		checkSt = checkers.WARNING
	}
	return checkers.NewChecker(checkSt, msg)
}
