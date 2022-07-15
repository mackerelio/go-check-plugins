package checkhttp

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
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
	URL                 string   `short:"u" long:"url" required:"true" description:"A URL to connect to"`
	Statuses            []string `short:"s" long:"status" description:"mapping of HTTP status"`
	NoCheckCertificate  bool     `long:"no-check-certificate" description:"Do not check certificate"`
	SourceIP            string   `short:"i" long:"source-ip" description:"source IP address"`
	Headers             []string `short:"H" description:"HTTP request headers"`
	Regexp              string   `short:"p" long:"pattern" description:"Expected pattern in the content"`
	MaxRedirects        int      `long:"max-redirects" description:"Maximum number of redirects followed" default:"10"`
	Method              string   `short:"m" long:"method" choice:"GET" choice:"HEAD" choice:"POST" choice:"PUT" default:"GET" description:"Specify a GET, HEAD, POST, or PUT operation"`
	ConnectTos          []string `long:"connect-to" value-name:"HOST1:PORT1:HOST2:PORT2" description:"Request to HOST2:PORT2 instead of HOST1:PORT1"`
	Proxy               string   `short:"x" long:"proxy" value-name:"[PROTOCOL://][USER:PASS@]HOST[:PORT]" description:"Use the specified proxy. PROTOCOL's default is http, and PORT's default is 1080."`
	BasicAuth           string   `long:"user" value-name:"USER[:PASSWORD]" description:"Basic Authentication user ID and an optional password."`
	RequireBytes        int64    `short:"B" long:"require-bytes" description:"Check the response contains exactly BYTES bytes" default:"-1"`
	Body                string   `short:"d" long:"body" description:"Send a data body string with the request"`
	MinBytes            int64    `short:"g" long:"min-bytes" description:"Check the response contains at least BYTES bytes" default:"-1"`
	Timeout             int64    `short:"t" long:"timeout" description:"Set the total execution timeout in seconds" default:"0"`
	CertFile            string   `long:"cert-file" description:"A Cert file to use for client authentication"`
	KeyFile             string   `long:"key-file" description:"A Key file to use for client authentication"`
	CaFile              string   `long:"ca-file" description:"A CA Cert file to use for client authentication"`
	NoResTimeSuccessMsg bool     `long:"no-restime-success-msg" description:"Do not output response time on success. Omissioning success report in mackerel-agent."`
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

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		InsecureSkipVerify: opts.NoCheckCertificate,
	}

	if opts.CertFile != "" && opts.KeyFile != "" {
		// Load client cert
		cert, err := tls.LoadX509KeyPair(opts.CertFile, opts.KeyFile)
		if err != nil {
			return checkers.Unknown(err.Error())
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
		tlsConfig.BuildNameToCertificate()
	}

	if opts.CaFile != "" {
		// Load CA cert
		caCert, err := ioutil.ReadFile(opts.CaFile)
		if err != nil {
			return checkers.Unknown(err.Error())
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
		tlsConfig.BuildNameToCertificate()
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
		Proxy:           http.ProxyFromEnvironment,
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
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(opts.Timeout),
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > opts.MaxRedirects {
			return http.ErrUseLastResponse
		}
		return nil
	}

	var req *http.Request
	if opts.Body == "" {
		req, err = http.NewRequest(opts.Method, opts.URL, nil)
		if err != nil {
			return checkers.Unknown(err.Error())
		}
	} else {
		body := strings.NewReader(opts.Body)
		req, err = http.NewRequest(opts.Method, opts.URL, body)
		if err != nil {
			return checkers.Unknown(err.Error())
		}
	}

	if len(opts.BasicAuth) != 0 {
		auth := strings.SplitN(opts.BasicAuth, ":", 2)
		if len(auth) == 1 {
			req.SetBasicAuth(auth[0], "")
		} else {
			req.SetBasicAuth(auth[0], auth[1])
		}
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

	if opts.RequireBytes != -1 && cLength != opts.RequireBytes {
		fmt.Fprintf(respMsg, "Response was '%d' bytes instead of '%d'\n", cLength, opts.RequireBytes)
		checkSt = checkers.CRITICAL
	}

	if opts.MinBytes != -1 && cLength < opts.MinBytes {
		fmt.Fprintf(respMsg, "Response was '%d' bytes instead of the indicated minimum '%d'\n", cLength, opts.MinBytes)
		checkSt = checkers.CRITICAL
	}

	if opts.NoResTimeSuccessMsg && checkSt == checkers.OK {
		fmt.Fprintf(respMsg, "%s %s - %d bytes",
			resp.Proto, resp.Status, cLength)
	} else {
		fmt.Fprintf(respMsg, "%s %s - %d bytes in %f second response time",
			resp.Proto, resp.Status, cLength, elapsed.Seconds())
	}

	return checkers.NewChecker(checkSt, respMsg.String())
}
