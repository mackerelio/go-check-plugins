package checkhttp

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
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
	ConnectTos         []string `long:"connect-to" value-name:"HOST1:PORT1:HOST2:PORT2" description:"Request to HOST2:PORT2 instead of HOST1:PORT1"`
	Proxy              string   `short:"x" long:"proxy" value-name:"[PROTOCOL://]HOST[:PORT]" description:"Use the specified proxy. PROTOCOL's default is http, and PORT's default is 1080."`
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

// when empty:
// - src* will be treated as ANY
// - dest* will be treated as unchanged
type resolveMapping struct {
	srcHost  string
	srcPort  string
	destHost string
	destPort string
}

func newReplacableDial(dialer *net.Dialer, mappings []resolveMapping) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, hostport string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(hostport)
		if err != nil {
			return nil, err
		}

		addr := hostport
		for _, m := range mappings {
			if m.srcHost != "" && m.srcHost != host {
				continue
			}
			if m.srcPort != "" && m.srcPort != port {
				continue
			}
			if m.destHost != "" {
				host = m.destHost
			}
			if m.destPort != "" {
				port = m.destPort
			}
			addr = net.JoinHostPort(host, port)
			break
		}
		return dialer.DialContext(ctx, network, addr)
	}
}

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

var connectToRegexp = regexp.MustCompile(`^(\[.+\]|[^\[\]]+)?:(\d*):(\[.+\]|[^\[\]]+)?:(\d+)?$`)

func parseConnectTo(opts *checkHTTPOpts) ([]resolveMapping, error) {
	mappings := make([]resolveMapping, len(opts.ConnectTos))
	for i, c := range opts.ConnectTos {
		s := connectToRegexp.FindStringSubmatch(c)
		if len(s) == 0 {
			return nil, fmt.Errorf("Invalid --connect-to pattern: %s", c)
		}
		r := resolveMapping{}
		if len(s) >= 2 {
			r.srcHost = s[1]
		}
		if len(s) >= 3 {
			r.srcPort = s[2]
		}
		if len(s) >= 4 {
			r.destHost = s[3]
		}
		if len(s) >= 5 {
			r.destPort = s[4]
		}
		mappings[i] = r
	}
	return mappings, nil
}

func parseProxy(opts *checkHTTPOpts) (*url.URL, error) {
	if opts.Proxy == "" {
		return nil, nil
	}
	// Append protocol if absent.
	// Overwriting u.Scheme is not enough, since url.Parse cannot parse "HOST:PORT" since it's ambiguous
	proxy := opts.Proxy
	if !strings.Contains(proxy, "://") {
		proxy = "http://" + proxy
	}
	u, err := url.Parse(proxy)
	if err != nil {
		return nil, err
	}
	if u.Port() == "" {
		u.Host = u.Hostname() + ":1080"
	}
	return u, nil
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
	// same as http.Transport's default dialer
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	if opts.SourceIP != "" {
		ip := net.ParseIP(opts.SourceIP)
		if ip == nil {
			return checkers.Unknown(fmt.Sprintf("Invalid source IP address: %v", opts.SourceIP))
		}
		dialer.LocalAddr = &net.TCPAddr{IP: ip}
	}

	proxyURL, err := parseProxy(&opts)
	if err != nil {
		return checkers.Unknown(err.Error())
	}
	if proxyURL != nil {
		tr.Proxy = http.ProxyURL(proxyURL)
	}

	if len(opts.ConnectTos) != 0 {
		resolves, err := parseConnectTo(&opts)
		if err != nil {
			return checkers.Unknown(err.Error())
		}
		tr.DialContext = newReplacableDial(dialer, resolves)
	}
	client := &http.Client{Transport: tr}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > opts.MaxRedirects {
			return http.ErrUseLastResponse
		}
		return nil
	}

	req, err := http.NewRequest(http.MethodGet, opts.URL, nil)
	if err != nil {
		return checkers.Unknown(err.Error())
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

	// set default User-Agent unless specified by `opts.Headers`
	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", "check-http")
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
