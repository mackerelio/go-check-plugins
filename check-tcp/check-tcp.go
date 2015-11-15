package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type tcpOpts struct {
	exchange
	Service  string  `long:"service"`
	Hostname string  `short:"H" long:"hostname" description:"Host name or IP Address"`
	Timeout  float64 `short:"t" long:"timeout" default:"10" description:"Seconds before connection times out"`
	MaxBytes int     `short:"m" long:"maxbytes"`
	Delay    float64 `short:"d" long:"delay" description:"Seconds to wait between sending string and polling for response"`
	Warning  float64 `short:"w" long:"warning" description:"Response time to result in warning status (seconds)"`
	Critical float64 `short:"c" long:"critical" description:"Response time to result in critical status (seconds)"`
	Escape   bool    `short:"E" long:"escape" description:"Can use \\n, \\r, \\t or \\ in send or quit string. Must come before send or quit option. By default, nothing added to send, \\r\\n added to end of quit"`
}

type exchange struct {
	Send          string `short:"s" long:"send" description:"String to send to the server"`
	ExpectPattern string `short:"e" long:"expect-pattern" description:"Regexp pattern to expect in server response"`
	Quit          string `short:"q" long:"quit" description:"String to send server to initiate a clean close of the connection"`
	Port          int    `short:"p" long:"port" description:"Port number"`
	SSL           bool   `short:"S" long:"ssl" description:"Use SSL for the connection."`
	expectReg     *regexp.Regexp
}

func main() {
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

func (opts *tcpOpts) prepare() error {
	opts.Service = strings.ToUpper(opts.Service)
	defaultEx := defaultExchange(opts.Service)
	opts.merge(defaultEx)

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

func defaultExchange(svc string) exchange {
	switch svc {
	case "FTP":
		return exchange{
			Port:          21,
			ExpectPattern: `^220`,
			Quit:          "QUIT",
		}
	case "POP":
		return exchange{
			Port:          110,
			ExpectPattern: `^\+OK`,
			Quit:          "QUIT",
		}
	case "SPOP":
		return exchange{
			Port:          995,
			ExpectPattern: `^\+OK`,
			Quit:          "QUIT",
			SSL:           true,
		}
	case "IMAP":
		return exchange{
			Port:          143,
			ExpectPattern: `^\* OK`,
			Quit:          "a1 LOGOUT",
		}
	case "SIMAP":
		return exchange{
			Port:          993,
			ExpectPattern: `^\* OK`,
			Quit:          "a1 LOGOUT",
			SSL:           true,
		}
	case "SMTP":
		return exchange{
			Port:          25,
			ExpectPattern: `^220`,
			Quit:          "QUIT",
		}
	case "SSMTP":
		return exchange{
			Port:          465,
			ExpectPattern: `^220`,
			Quit:          "QUIT",
			SSL:           true,
		}
	}
	return exchange{}
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
}

func dial(address string, ssl bool) (net.Conn, error) {
	if ssl {
		return tls.Dial("tcp", address, &tls.Config{})
	}
	return net.Dial("tcp", address)
}

func (opts *tcpOpts) run() *checkers.Checker {
	err := opts.prepare()
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	send := opts.Send
	expect := opts.ExpectPattern
	quit := opts.Quit
	address := fmt.Sprintf("%s:%d", opts.Hostname, opts.Port)

	start := time.Now()
	if opts.Delay > 0 {
		time.Sleep(time.Duration(opts.Delay) * time.Second)
	}
	conn, err := dial(address, opts.SSL)
	if err != nil {
		return checkers.Critical(err.Error())
	}
	defer conn.Close()

	if send != "" {
		err := write(conn, []byte(send), opts.Timeout)
		if err != nil {
			return checkers.Critical(err.Error())
		}
	}

	res := ""
	if opts.ExpectPattern != "" {
		buf, err := slurp(conn, opts.MaxBytes, opts.Timeout)
		if err != nil {
			return checkers.Critical(err.Error())
		}

		res = string(buf)
		if expect != "" && !opts.expectReg.MatchString(res) {
			return checkers.Critical("Unexpected response from host/socket: " + res)
		}
	}

	if quit != "" {
		err := write(conn, []byte(quit), opts.Timeout)
		if err != nil {
			return checkers.Critical(err.Error())
		}
	}
	elapsed := time.Now().Sub(start)

	chkSt := checkers.OK
	if opts.Warning > 0 && elapsed > time.Duration(opts.Warning)*time.Second {
		chkSt = checkers.WARNING
	}
	if opts.Critical > 0 && elapsed > time.Duration(opts.Critical)*time.Second {
		chkSt = checkers.CRITICAL
	}
	msg := fmt.Sprintf("%.3f seconds response time on", float64(elapsed)/float64(time.Second))
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

func write(conn net.Conn, content []byte, timeout float64) error {
	if timeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
	_, err := conn.Write(content)
	return err
}

func slurp(conn net.Conn, maxbytes int, timeout float64) ([]byte, error) {
	buf := []byte{}
	readLimit := 32 * 1024
	if maxbytes > 0 {
		readLimit = maxbytes
	}
	readBytes := 0
	if timeout > 0 {
		conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
	for {
		tmpBuf := make([]byte, readLimit)
		i, err := conn.Read(tmpBuf)
		if err != nil {
			return buf, err
		}
		buf = append(buf, tmpBuf[:i]...)
		readBytes += i
		if i < readLimit || (maxbytes > 0 && maxbytes <= readBytes) {
			break
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
