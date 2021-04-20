package checksmtp

import (
	"net"
	"testing"
	"time"

	"net/textproto"
	"strings"
	"sync/atomic"
)

var responseDefault = []string{
	"220 mail.example.com ESMTP unknown",
	"250 mail.example.com",
	"250 OK",
	"250 OK",
	"221 Bye",
}

type mockSMTPServer struct {
	listener  net.Listener
	delay     int
	responses []string
	closed    int32
}

func (m *mockSMTPServer) runServe() error {
	doneCh := make(chan struct{})
	errCh := make(chan error)
	go func() {
		time.Sleep(time.Duration(m.delay) * time.Second)

		conn, err := m.listener.Accept()
		if err != nil {
			errCh <- err
			return
		}
		tc := textproto.NewConn(conn)

		for _, res := range m.responses {
			if err := tc.PrintfLine(res); err != nil {
				errCh <- err
				return
			}
		}
		doneCh <- struct{}{}
	}()

	select {
	case <-doneCh:
		return nil
	case err := <-errCh:
		if atomic.LoadInt32(&m.closed) != 0 {
			return nil
		}
		return err
	}
}

func (m *mockSMTPServer) Close() {
	atomic.StoreInt32(&m.closed, 1)
	m.listener.Close()
}

func TestSMTP_Default(t *testing.T) {
	var cases = []struct {
		name      string
		argv      []string
		delay     int
		expected  string
		responses []string
	}{
		{
			name:      "warning",
			argv:      []string{"--warning", "1"},
			delay:     2,
			responses: responseDefault,
			expected:  "WARNING",
		},
		{
			name:      "critical",
			argv:      []string{"--critical", "1"},
			delay:     2,
			responses: responseDefault,
			expected:  "CRITICAL",
		},
		{
			name:      "timeout",
			argv:      []string{"--critical", "1", "--timeout", "-1"},
			delay:     1,
			responses: responseDefault,
			expected:  "CRITICAL",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ln, err := net.Listen("tcp", "127.0.0.1:0")
			defer ln.Close()
			if err != nil {
				t.Fatal(err.Error())
			}
			s := mockSMTPServer{
				listener:  ln,
				responses: c.responses,
				delay:     c.delay,
			}
			errCh := make(chan error, 1)
			go func() {
				if err := s.runServe(); err != nil {
					errCh <- err
				}
				close(errCh)
			}()
			port := strings.Split(ln.Addr().String(), ":")[1]
			argv := []string{"--host", "127.0.0.1", "--port", port}
			argv = append(argv, c.argv...)
			ckr := run(argv)
			s.Close()
			if err := <-errCh; err != nil {
				t.Fatal(err)
			}
			if ckr.Status.String() != c.expected {
				t.Errorf("%s: %s", ckr.Status.String(), ckr.Message)
			}
		})
	}
}

func TestSMTP_NoThresholdOptions(t *testing.T) {
	ckr := run([]string{})
	if ckr.Status.String() != "UNKNOWN" {
		t.Errorf("expected UNKNOWN to eq %s", ckr.Status.String())
	}
}
