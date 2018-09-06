package checkhttp

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

// XXX more options
type checkHTTPOpts struct {
	URL                string   `short:"u" long:"url" required:"true" description:"A URL to connect to"`
	Statuses           []string `short:"s" long:"status" description:"mapping of HTTP status"`
	NoCheckCertificate bool     `long:"no-check-certificate" description:"Do not check certificate"`
	SourceIP           string   `short:"i" long:"source-ip" description:"source IP address"`
	Headers            []string `short:"H" description:"HTTP request headers"`
	Regexp             string   `short:"p" long:"pattern" description:"Expected pattern in the content"`
	MaxRedirects       int      `long:"max-redirects" description:"Maximum number of redirects followed" default:"10"`
}

// Do the plugin
func Do() {
	ckr := Run(os.Args[1:])
	ckr.Name = "HTTP"
	ckr.Exit()
}

type statusRange struct {
	min     int
	max     int
	checkSt checkers.Status
}

const invalidMapping = "Invalid mapping of status: %s"

func parseStatusRanges(opts *checkHTTPOpts) ([]statusRange, error) {
	var statuses []statusRange
	for _, s := range opts.Statuses {
		token := strings.SplitN(s, "=", 2)
		if len(token) != 2 {
			return nil, fmt.Errorf(invalidMapping, s)
		}
		values := strings.Split(token[0], "-")

		var r statusRange
		var err error

		switch len(values) {
		case 1:
			r.min, err = strconv.Atoi(values[0])
			if err != nil {
				return nil, fmt.Errorf(invalidMapping, s)
			}
			r.max = r.min
		case 2:
			r.min, err = strconv.Atoi(values[0])
			if err != nil {
				return nil, fmt.Errorf(invalidMapping, s)
			}
			r.max, err = strconv.Atoi(values[1])
			if err != nil {
				return nil, fmt.Errorf(invalidMapping, s)
			}
			if r.min > r.max {
				return nil, fmt.Errorf(invalidMapping, s)
			}
		default:
			return nil, fmt.Errorf(invalidMapping, s)
		}

		switch strings.ToUpper(token[1]) {
		case "OK":
			r.checkSt = checkers.OK
		case "WARNING":
			r.checkSt = checkers.WARNING
		case "CRITICAL":
			r.checkSt = checkers.CRITICAL
		case "UNKNOWN":
			r.checkSt = checkers.UNKNOWN
		default:
			return nil, fmt.Errorf(invalidMapping, s)
		}
		statuses = append(statuses, r)
	}
	return statuses, nil
}

func parseHeader(opts *checkHTTPOpts) (http.Header, error) {
	reader := bufio.NewReader(strings.NewReader(strings.Join(opts.Headers, "\r\n") + "\r\n\r\n"))
	tp := textproto.NewReader(reader)
	mimeheader, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, fmt.Errorf("failed to parse header: %s", err)
	}
	return http.Header(mimeheader), nil
}

// Run do external monitoring via HTTP
func Run(args []string) *checkers.Checker {
	opts := checkHTTPOpts{}
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	statusRanges, err := parseStatusRanges(&opts)
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.NoCheckCertificate,
		},
		Proxy: http.ProxyFromEnvironment,
	}
	if opts.SourceIP != "" {
		ip := net.ParseIP(opts.SourceIP)
		if ip == nil {
			return checkers.Unknown(fmt.Sprintf("Invalid source IP address: %v", opts.SourceIP))
		}
		tr.Dial = (&net.Dialer{
			LocalAddr: &net.TCPAddr{IP: ip},
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).Dial
	}
	client := &http.Client{Transport: tr}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > opts.MaxRedirects {
			return http.ErrUseLastResponse
		}
		return nil
	}

	req, err := http.NewRequest("GET", opts.URL, nil)
	if err != nil {
		return checkers.Critical(err.Error())
	}

	if len(opts.Headers) != 0 {
		header, err := parseHeader(&opts)
		if err != nil {
			return checkers.Unknown(err.Error())
		}

		// Host header must be set via req.Host
		if host := header.Get("Host"); len(host) != 0 {
			req.Host = host
			header.Del("Host")
		}

		req.Header = header
	}

	stTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return checkers.Critical(err.Error())
	}
	elapsed := time.Since(stTime)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	cLength := resp.ContentLength
	if cLength == -1 {
		cLength = int64(len(body))
	}

	checkSt := checkers.UNKNOWN

	found := false
	for _, st := range statusRanges {
		if st.min <= resp.StatusCode && resp.StatusCode <= st.max {
			checkSt = st.checkSt
			found = true
			break
		}
	}

	if !found {
		switch st := resp.StatusCode; true {
		case st < 400:
			checkSt = checkers.OK
		case st < 500:
			checkSt = checkers.WARNING
		default:
			checkSt = checkers.CRITICAL
		}
	}

	respMsg := new(bytes.Buffer)

	if opts.Regexp != "" {
		re, err := regexp.Compile(opts.Regexp)
		if err != nil {
			return checkers.Unknown(err.Error())
		}
		if !re.Match(body) {
			fmt.Fprintf(respMsg, "'%s' not found in the content\n", opts.Regexp)
			checkSt = checkers.CRITICAL
		}
	}

	fmt.Fprintf(respMsg, "%s %s - %d bytes in %f second response time",
		resp.Proto, resp.Status, cLength, elapsed.Seconds())

	return checkers.NewChecker(checkSt, respMsg.String())
}
