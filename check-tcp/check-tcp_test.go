package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestEscapedString(t *testing.T) {
	assert.Equal(t, "\n", escapedString(`\n`), "something went wrong")
	assert.Equal(t, "hoge\\", escapedString(`hoge\`), "something went wrong")
	assert.Equal(t, "ho\rge", escapedString(`ho\rge`), "something went wrong")
	assert.Equal(t, "ho\\oge", escapedString(`ho\oge`), "something went wrong")
	assert.Equal(t, "", escapedString(``), "something went wrong")
}

func TestTLS(t *testing.T) {
	opts, err := parseArgs([]string{"-S", "-H", "www.verisign.com", "-p", "443"})
	assert.Equal(t, nil, err, "no errors")
	ckr := opts.run()
	assert.Equal(t, checkers.OK, ckr.Status, "should be OK")
}

func TestFTP(t *testing.T) {
	opts, err := parseArgs([]string{"--service=ftp", "-H", "ftp.iij.ad.jp"})
	assert.Equal(t, nil, err, "no errors")
	ckr := opts.run()
	assert.Equal(t, checkers.OK, ckr.Status, "should be OK")
}

func TestHTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(time.Second / 5)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "OKOK")
	}))
	defer ts.Close()

	u, _ := url.Parse(ts.URL)

	host, port, _ := net.SplitHostPort(u.Host)

	testOk := func() {
		opts, err := parseArgs([]string{"-H", host, "-p", port, "--send", `GET / HTTP/1.1\r\n\r\n`, "-E", "-e", "OKOK"})
		assert.Equal(t, nil, err, "no errors")
		ckr := opts.run()
		assert.Equal(t, checkers.OK, ckr.Status, "should be OK")
		assert.Regexp(t, `seconds response time on`, ckr.Message, "Unexpected response")
	}
	testOk()

	testUnexpected := func() {
		opts, err := parseArgs(
			[]string{"-H", host, "-p", port, "--send", `GET / HTTP/1.1\r\n\r\n`, "-E", "-e", "OKOKOK"})
		assert.Equal(t, nil, err, "no errors")
		ckr := opts.run()
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "should be CRITICAL")
		assert.Regexp(t, `Unexpected response from`, ckr.Message, "Unexpected response")
	}
	testUnexpected()

	testOverWarn := func() {
		opts, err := parseArgs(
			[]string{"-H", host, "-p", port, "--send", `GET / HTTP/1.1\r\n\r\n`, "-E", "-e", "OKOK", "-w", "0.1"})
		assert.Equal(t, nil, err, "no errors")
		ckr := opts.run()
		assert.Equal(t, checkers.WARNING, ckr.Status, "should be Warning")
		assert.Regexp(t, `seconds response time on`, ckr.Message, "Unexpected response")
	}
	testOverWarn()

	testOverCrit := func() {
		opts, err := parseArgs(
			[]string{"-H", host, "-p", port, "--send", "GET / HTTP/1.1\r\n\r\n", "-e", "OKOK", "-c", "0.1"})
		assert.Equal(t, nil, err, "no errors")
		ckr := opts.run()
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "should be Critical")
		assert.Regexp(t, `seconds response time on`, ckr.Message, "Unexpected response")
	}
	testOverCrit()
}

func TestUnixDomainSocket(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	sock := fmt.Sprintf("%s/test.sock", dir)

	l, err := net.Listen("unix", sock)
	if err != nil {
		t.Error(err)
	}

	go func() {
		for {
			ls, err := l.Accept()
			if err != nil {
				t.Error(err)
			}
			go func(c net.Conn) {
				defer c.Close()

				buf := make([]byte, 1024)

				_, err := c.Read(buf)

				if err == io.EOF {
					return
				}

				c.Write([]byte("OKOK"))
			}(ls)
		}
	}()

	testOk := func() {
		opts, err := parseArgs([]string{"-U", sock, "--send", `PING`, "-E", "-e", "OKOK"})
		assert.Equal(t, nil, err, "no errors")
		ckr := opts.run()
		assert.Equal(t, checkers.OK, ckr.Status, "should be OK")
		assert.Regexp(t, `seconds response time on`, ckr.Message, "Unexpected response")
	}
	testOk()

	testUnexpected := func() {
		opts, err := parseArgs([]string{"-U", sock, "--send", `PING`, "-E", "-e", "OKOKOK"})
		assert.Equal(t, nil, err, "no errors")
		ckr := opts.run()
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "should be CRITICAL")
		assert.Regexp(t, `Unexpected response from`, ckr.Message, "Unexpected response")
	}
	testUnexpected()

	testOverWarn := func() {
		opts, err := parseArgs([]string{"-U", sock, "--send", `PING`, "-E", "-e", "OKOK", "-w", "0.1"})
		assert.Equal(t, nil, err, "no errors")
		ckr := opts.run()
		assert.Equal(t, checkers.WARNING, ckr.Status, "should be Warning")
		assert.Regexp(t, `seconds response time on`, ckr.Message, "Unexpected response")
	}
	testOverWarn()

	testOverCrit := func() {
		opts, err := parseArgs([]string{"-U", sock, "--send", `PING`, "-E", "-e", "OKOK", "-c", "0.1"})
		assert.Equal(t, nil, err, "no errors")
		ckr := opts.run()
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "should be Critical")
		assert.Regexp(t, `seconds response time on`, ckr.Message, "Unexpected response")
	}
	testOverCrit()
}

func TestHTTPIPv6(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(time.Second / 5)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "OKOK")
	})

	l, err := net.Listen("tcp", "[::1]:0")
	if err != nil {
		t.Error(err)
	}
	defer l.Close()
	h, port, _ := net.SplitHostPort(l.Addr().String())
	host := fmt.Sprintf("[%s]", h)

	go func() {
		for {
			http.Serve(l, nil)
		}
	}()

	testOk := func() {
		opts, err := parseArgs([]string{"-H", host, "-p", port, "--send", `GET / HTTP/1.1\r\n\r\n`, "-E", "-e", "OKOK"})
		assert.Equal(t, nil, err, "no errors")
		ckr := opts.run()
		assert.Equal(t, checkers.OK, ckr.Status, "should be OK")
		assert.Regexp(t, `seconds response time on`, ckr.Message, "Unexpected response")
	}
	testOk()

	testUnexpected := func() {
		opts, err := parseArgs(
			[]string{"-H", host, "-p", port, "--send", `GET / HTTP/1.1\r\n\r\n`, "-E", "-e", "OKOKOK"})
		assert.Equal(t, nil, err, "no errors")
		ckr := opts.run()
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "should be CRITICAL")
		assert.Regexp(t, `Unexpected response from`, ckr.Message, "Unexpected response")
	}
	testUnexpected()

	testOverWarn := func() {
		opts, err := parseArgs(
			[]string{"-H", host, "-p", port, "--send", `GET / HTTP/1.1\r\n\r\n`, "-E", "-e", "OKOK", "-w", "0.1"})
		assert.Equal(t, nil, err, "no errors")
		ckr := opts.run()
		assert.Equal(t, checkers.WARNING, ckr.Status, "should be Warning")
		assert.Regexp(t, `seconds response time on`, ckr.Message, "Unexpected response")
	}
	testOverWarn()

	testOverCrit := func() {
		opts, err := parseArgs(
			[]string{"-H", host, "-p", port, "--send", "GET / HTTP/1.1\r\n\r\n", "-e", "OKOK", "-c", "0.1"})
		assert.Equal(t, nil, err, "no errors")
		ckr := opts.run()
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "should be Critical")
		assert.Regexp(t, `seconds response time on`, ckr.Message, "Unexpected response")
	}
	testOverCrit()
}
