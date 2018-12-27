package checkjson

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	ckr := Run([]string{"-u", "hoge"})
	assert.Equal(t, ckr.Status, checkers.CRITICAL, "chr.Status should be CRITICAL")
	assert.Equal(t, ckr.Message, `Get hoge: unsupported protocol scheme ""`, "something went wrong")
}

func TestNoCheckCertificate(t *testing.T) {
	ckr := Run([]string{"--no-check-certificate", "-u", "hoge"})
	assert.Equal(t, ckr.Status, checkers.CRITICAL, "chr.Status should be CRITICAL")
	assert.Equal(t, ckr.Message, `Get hoge: unsupported protocol scheme ""`, "something went wrong")
}

func TestStatusRange(t *testing.T) {
	tests := []struct {
		args []string
		want checkers.Status
		err  bool
	}{
		{
			args: []string{"-s", "404=ok", "-u", "hoge"},
			want: checkers.CRITICAL,
		},
		{
			args: []string{"-s", "404=ok", "-u", "https://mackerel.io/404"},
			want: checkers.OK,
		},
		{
			args: []string{"-s", "401=ok", "-u", "https://mackerel.io/404"},
			want: checkers.WARNING,
		},
		{
			args: []string{"-s", "300-404=ok", "-u", "https://mackerel.io/404"},
			want: checkers.OK,
		},
		{
			args: []string{"-s", "200-300-404=ok", "-u", "http://example.com"},
			want: checkers.UNKNOWN,
		},
		{
			args: []string{"-s", "300=ok", "-s", "404=ok", "-u", "https://mackerel.io/404"},
			want: checkers.OK,
		},
		{
			args: []string{"-s", "300", "-s", "404=ok", "-u", "https://mackerel.io/404"},
			want: checkers.UNKNOWN,
		},
		{
			args: []string{"-s", "=ok", "-s", "404=ok", "-u", "https://mackerel.io/404"},
			want: checkers.UNKNOWN,
		},
		{
			args: []string{"-s", "=ok", "-s", "404-200=ok", "-u", "https://mackerel.io/404"},
			want: checkers.UNKNOWN,
		},
	}
	for _, tt := range tests {
		ckr := Run(tt.args)
		assert.Equal(t, ckr.Status, tt.want, fmt.Sprintf("chr.Status wrong: %v", ckr.Status))
	}
}

func TestSourceIP(t *testing.T) {
	ckr := Run([]string{"-u", "hoge", "-i", "1.2.3"})
	assert.Equal(t, ckr.Status, checkers.UNKNOWN, "chr.Status should be UNKNOWN")
}

func TestHost(t *testing.T) {
	testHost := "mackerel.io"
	testHeader := "test"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, testHost, r.Host)
		header := r.Header
		assert.Equal(t, testHeader, header.Get("TestHeader"))
	}))
	defer ts.Close()

	ckr := Run([]string{
		"-H", fmt.Sprintf("Host: %s", testHost),
		"-H", fmt.Sprintf("TestHeader: %s", testHeader),
		"-u", ts.URL,
	})

	assert.Equal(t, ckr.Status, checkers.OK, "ckr.Status should be OK")
}

func TestExpectedContent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	testCases := []struct {
		regexp string
		status checkers.Status
	}{
		{
			regexp: "Hello, client",
			status: checkers.OK,
		},
		{
			regexp: "Wrong response",
			status: checkers.CRITICAL,
		},
		{
			regexp: "Hel.*",
			status: checkers.OK,
		},
		{
			regexp: "clientt?",
			status: checkers.OK,
		},
		{
			regexp: "???",
			status: checkers.UNKNOWN,
		},
	}

	for i, tc := range testCases {
		ckr := Run([]string{"-u", ts.URL, "-p", tc.regexp})
		assert.Equal(t, ckr.Status, tc.status, "#%d: Status should be %s", i, tc.status)
	}
}

func TestMaxRedirects(t *testing.T) {
	redirectedPath := "/redirected"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != redirectedPath {
			http.Redirect(w, r, redirectedPath, 301)
		}
	}))
	defer ts.Close()

	testCases := []struct {
		args []string
		want checkers.Status
	}{
		{
			args: []string{"-s", "301=ok", "-s", "100-300=warning", "-s", "302-599=warning",
				"-u", ts.URL, "--max-redirects", "0"},
			want: checkers.OK,
		},
		{
			args: []string{"-s", "200=ok", "-s", "100-199=warning", "-s", "201-599=warning",
				"-u", ts.URL, "--max-redirects", "1"},
			want: checkers.OK,
		},
		{
			args: []string{"-s", "200=ok", "-s", "100-199=warning", "-s", "201-599=warning",
				"-u", ts.URL},
			want: checkers.OK,
		},
	}

	for i, tc := range testCases {
		ckr := Run(tc.args)
		assert.Equal(t, ckr.Status, tc.want, "#%d: Status should be %s", i, tc.want)
	}
}

func TestConnectTos(t *testing.T) {
	// expected server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "this is %s\n", r.URL.Host)
	}))
	defer ts.Close()
	// extract host and port
	s := strings.SplitN(ts.URL, ":", 3)
	addr := strings.TrimPrefix(s[1], "//")
	port := s[2]

	// NON-expected server
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "wrong server!", 500)
	}))
	defer ts2.Close()
	// extract host and port
	s2 := strings.SplitN(ts.URL, ":", 3)
	addr2 := strings.TrimPrefix(s2[1], "//")
	port2 := s2[2]

	testCases := []struct {
		args []string
		want checkers.Status
	}{
		{
			// not affected at all
			args: []string{"--connect-to", fmt.Sprintf("hoge:80:%s:%s", addr, port),
				"-u", ts.URL},
			want: checkers.OK,
		},
		{
			// connected to target
			args: []string{"--connect-to", fmt.Sprintf("hoge:80:%s:%s", addr, port),
				"-u", "http://hoge"},
			want: checkers.OK,
		},
		{
			// empty srcHost means ANY
			args: []string{"--connect-to", fmt.Sprintf(":80:%s:%s", addr, port),
				"-u", "http://hoge"},
			want: checkers.OK,
		},
		{
			// empty srcPort means ANY
			args: []string{"--connect-to", fmt.Sprintf("hoge::%s:%s", addr, port),
				"-u", "http://hoge"},
			want: checkers.OK,
		},
		{
			// empty destHost means unchanged
			args: []string{"--connect-to", fmt.Sprintf("%s:80::%s", addr, port),
				"-u", fmt.Sprintf("http://%s", addr)},
			want: checkers.OK,
		},
		{
			// empty destPort means unchanged
			args: []string{"--connect-to", fmt.Sprintf("hoge:%s:%s:", port, addr),
				"-u", fmt.Sprintf("http://hoge:%s", port)},
			want: checkers.OK,
		},
		{
			// host mismatch ignored
			args: []string{"--connect-to", fmt.Sprintf("not.target:%s:%s:%s", port, addr2, port2),
				"-u", ts.URL},
			want: checkers.OK,
		},
		{
			// port mismatch ignored
			args: []string{"--connect-to", fmt.Sprintf("%s:%s:%s:%s", addr, port2, addr2, port2),
				"-u", ts.URL},
			want: checkers.OK,
		},
		{
			// host mismatch ignored, even if port is empty
			args: []string{"--connect-to", fmt.Sprintf("not.target::%s:%s", addr2, port2),
				"-u", ts.URL},
			want: checkers.OK,
		},
		{
			// port mismatch ignored, even if host is empty
			args: []string{"--connect-to", fmt.Sprintf(":%s:%s:%s", port2, addr2, port2),
				"-u", ts.URL},
			want: checkers.OK,
		},
		{
			// multiple setting (1)
			args: []string{"--connect-to", fmt.Sprintf("not.hoge:80:%s:%s", addr2, port2),
				"--connect-to", fmt.Sprintf("hoge:80:%s:%s", addr, port),
				"-u", "http://hoge"},
			want: checkers.OK,
		},
		{
			// multiple setting (2)
			args: []string{"--connect-to", fmt.Sprintf("hoge:80:%s:%s", addr, port),
				"--connect-to", fmt.Sprintf("hoge:80:%s:%s", addr2, port2),
				"-u", "http://hoge"},
			want: checkers.OK,
		},
		{
			// Invalid pattern
			args: []string{"--connect-to", "foo:123:",
				"-u", ts.URL},
			want: checkers.UNKNOWN,
		},
	}

	for i, tc := range testCases {
		ckr := Run(tc.args)
		assert.Equal(t, ckr.Status, tc.want, "#%d: Status should be %s, %s", i, tc.want, ckr.Message)
	}
}
