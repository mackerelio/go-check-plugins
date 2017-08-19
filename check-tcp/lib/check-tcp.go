package checktcp

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type tcpOpts struct {
	Service  string `long:"service" description:"Service name. e.g. ftp, smtp, pop, imap and so on"`
	Hostname string `short:"H" long:"hostname" description:"Host name or IP Address"`
	exchange
	Timeout  float64 `short:"t" long:"timeout" default:"10" description:"Seconds before connection times out"`
	MaxBytes int     `short:"m" long:"maxbytes" description:"Close connection once more than this number of bytes are received"`
	Delay    float64 `short:"d" long:"delay" description:"Seconds to wait between sending string and polling for response"`
	Warning  float64 `short:"w" long:"warning" description:"Response time to result in warning status (seconds)"`
	Critical float64 `short:"c" long:"critical" description:"Response time to result in critical status (seconds)"`
	Escape   bool    `short:"E" long:"escape" description:"Can use \\n, \\r, \\t or \\ in send or quit string. Must come before send or quit option. By default, nothing added to send, \\r\\n added to end of quit"`
}

type exchange struct {
	Port               int    `short:"p" long:"port" description:"Port number"`
	Send               string `short:"s" long:"send" description:"String to send to the server"`
	ExpectPattern      string `short:"e" long:"expect-pattern" description:"Regexp pattern to expect in server response"`
	Quit               string `short:"q" long:"quit" description:"String to send server to initiate a clean close of the connection"`
	SSL                bool   `short:"S" long:"ssl" description:"Use SSL for the connection."`
	UnixSock           string `short:"U" long:"unix-sock" description:"Unix Domain Socket"`
	NoCheckCertificate bool   `long:"no-check-certificate" description:"Do not check certificate"`
	expectReg          *regexp.Regexp
}

// Do the plugin
func Do() {
	opts, err := parseArgs(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
	ckr := opts.run()
	ckr.Name = "TCP"
	if opts.Service != "" {
		ckr.Name = opts.Service
	}
	ckr.Exit()
}

func parseArgs(args []string) (*tcpOpts, error) {
	opts := &tcpOpts{}
	_, err := flags.ParseArgs(opts, args)
	return opts, err
}

var defaultExchangeMap = map[string]exchange{
	"FTP": {
		Port:          21,
		ExpectPattern: `^220`,
		Quit:          "QUIT",
	},
	"POP": {
		Port:          110,
		ExpectPattern: `^\+OK`,
		Quit:          "QUIT",
	},
	"SPOP": {
		Port:          995,
		ExpectPattern: `^\+OK`,
		Quit:          "QUIT",
		SSL:           true,
	},
	"IMAP": {
		Port:          143,
		ExpectPattern: `^\* OK`,
		Quit:          "a1 LOGOUT",
	},
	"SIMAP": {
		Port:          993,
		ExpectPattern: `^\* OK`,
		Quit:          "a1 LOGOUT",
		SSL:           true,
	},
	"SMTP": {
		Port:          25,
		ExpectPattern: `^220`,
		Quit:          "QUIT",
	},
	"SSMTP": {
		Port:          465,
		ExpectPattern: `^220`,
		Quit:          "QUIT",
		SSL:           true,
	},
	"GEARMAN": {
		Port:          7003,
		Send:          "version\n",
		ExpectPattern: `\A[0-9]+\.[0-9]+\n\z`,
	},
}

func (opts *tcpOpts) prepare() error {
	opts.Service = strings.ToUpper(opts.Service)

	if opts.Service != "" {
		defaultEx, ok := defaultExchangeMap[opts.Service]
		if !ok {
			return fmt.Errorf("check-tcp called with unknown service: %s", opts.Service)
		}
		opts.merge(defaultEx)
	}

	if opts.Escape {
		opts.Quit = escapedString(opts.Quit)
		opts.Send = escapedString(opts.Send)
	} else if opts.Quit != "" {
		opts.Quit += "\r\n"
	}
	var err error
	if opts.ExpectPattern != "" {
		opts.expectReg, err = regexp.Compile(opts.ExpectPattern)
	}
	return err
}

func (opts *tcpOpts) merge(ex exchange) {
	if opts.Port == 0 {
		opts.Port = ex.Port
	}
	if opts.Send == "" {
		opts.Send = ex.Send
	}
	if opts.ExpectPattern == "" {
		opts.ExpectPattern = ex.ExpectPattern
	}
	if opts.Quit == "" {
		opts.Quit = ex.Quit
	}
	if !opts.SSL {
		opts.SSL = ex.SSL
	}
}

func dial(network, address string, ssl bool, noCheckCertificate bool, timeout time.Duration) (net.Conn, error) {
	d := &net.Dialer{Timeout: timeout}
	if ssl {
		return tls.DialWithDialer(d, network, address, &tls.Config{
			InsecureSkipVerify: noCheckCertificate,
		})
	}
	return d.Dial(network, address)
}

func (opts *tcpOpts) run() *checkers.Checker {
	err := opts.prepare()
	if err != nil {
		return checkers.Unknown(err.Error())
	}
	// prevent changing output of some commands
	os.Setenv("LANG", "C")
	os.Setenv("LC_ALL", "C")

	proto := "tcp"
	addr := fmt.Sprintf("%s:%d", opts.Hostname, opts.Port)
	if opts.UnixSock != "" {
		proto = "unix"
		addr = opts.UnixSock
	}
	timeout := time.Duration(opts.Timeout * float64(time.Second))
	start := time.Now()
	if opts.Delay > 0 {
		time.Sleep(time.Duration(opts.Delay) * time.Second)
	}

	conn, err := dial(proto, addr, opts.SSL, opts.NoCheckCertificate, timeout)
	if err != nil {
		return checkers.Critical(err.Error())
	}
	defer conn.Close()

	if opts.Send != "" {
		err := write(conn, []byte(opts.Send), timeout)
		if err != nil {
			return checkers.Critical(err.Error())
		}
	}

	res := ""
	if opts.expectReg != nil {
		buf, err := slurp(conn, opts.MaxBytes, timeout)
		if err != nil {
			return checkers.Critical(err.Error())
		}
		res = string(buf)
		if !opts.expectReg.MatchString(res) {
			return checkers.Critical("Unexpected response from host/socket: " + res)
		}
	}

	if opts.Quit != "" {
		err := write(conn, []byte(opts.Quit), timeout)
		if err != nil {
			return checkers.Critical(err.Error())
		}
	}
	elapsedSeconds := float64(time.Now().Sub(start)) / float64(time.Second)

	chkSt := checkers.OK
	if opts.Warning > 0 && elapsedSeconds > opts.Warning {
		chkSt = checkers.WARNING
	}
	if opts.Critical > 0 && elapsedSeconds > opts.Critical {
		chkSt = checkers.CRITICAL
	}
	msg := fmt.Sprintf("%.3f seconds response time on", elapsedSeconds)
	if opts.Hostname != "" {
		msg += " " + opts.Hostname
	}
	if opts.Port > 0 {
		msg += fmt.Sprintf(" port %d", opts.Port)
	}
	if res != "" {
		msg += fmt.Sprintf(" [%s]", strings.Trim(res, "\r\n"))
	}
	return checkers.NewChecker(chkSt, msg)
}

func write(conn net.Conn, content []byte, timeout time.Duration) error {
	if timeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(timeout))
	}
	_, err := conn.Write(content)
	return err
}

func slurp(conn net.Conn, maxbytes int, timeout time.Duration) ([]byte, error) {
	buf := []byte{}
	readLimit := 32 * 1024
	if maxbytes > 0 {
		readLimit = maxbytes
	}
	readBytes := 0
	if timeout > 0 {
		conn.SetReadDeadline(time.Now().Add(timeout))
	}
	for {
		tmpBuf := make([]byte, readLimit)
		i, err := conn.Read(tmpBuf)
		if i > 0 {
			buf = append(buf, tmpBuf[:i]...)
			readBytes += i
			if i < readLimit || (maxbytes > 0 && maxbytes <= readBytes) {
				break
			}
		}
		if err == io.EOF {
			return buf, nil
		}
		if err != nil {
			return buf, err
		}
	}
	return buf, nil
}

func escapedString(str string) (escaped string) {
	l := len(str)
	for i := 0; i < l; i++ {
		c := str[i]
		if c == '\\' && i+1 < l {
			i++
			c := str[i]
			switch c {
			case 'n':
				escaped += "\n"
			case 'r':
				escaped += "\r"
			case 't':
				escaped += "\t"
			case '\\':
				escaped += `\`
			default:
				escaped += `\` + string(c)
			}
		} else {
			escaped += string(c)
		}
	}
	return escaped
}
