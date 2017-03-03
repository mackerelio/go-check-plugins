package checkhttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"crypto/tls"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

// XXX more options
var opts struct {
	URL string `short:"u" long:"url" required:"true" description:"A URL to connect to"`
	NoCheckCertificate bool `long:"no-check-certificate" description:"Do not check certificate"`
}

// Do the plugin
func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "HTTP"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	tr := &http.Transport {
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.NoCheckCertificate,
		},
	}
	client := &http.Client{ Transport: tr }

	stTime := time.Now()
	resp, err := client.Get(opts.URL)
	if err != nil {
		return checkers.Critical(err.Error())
	}
	elapsed := time.Since(stTime)
	defer resp.Body.Close()

	cLength := resp.ContentLength
	if cLength == -1 {
		byt, _ := ioutil.ReadAll(resp.Body)
		cLength = int64(len(byt))
	}

	checkSt := checkers.UNKNOWN
	switch st := resp.StatusCode; true {
	case st < 400:
		checkSt = checkers.OK
	case st < 500:
		checkSt = checkers.WARNING
	default:
		checkSt = checkers.CRITICAL
	}

	msg := fmt.Sprintf("%s %s - %d bytes in %f second respons time",
		resp.Proto, resp.Status, cLength, elapsed.Seconds())

	return checkers.NewChecker(checkSt, msg)
}
