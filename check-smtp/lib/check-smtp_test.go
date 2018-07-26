package checksmtp

import (
	"net"
	"testing"
	"time"

	"net/textproto"
	"strings"
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
		return err
	}
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
			name:      "default",
			argv:      []string{},
			responses: responseDefault,
			expected:  "OK",
		},
		{
			name:      "warning",
			argv:      []string{"--warning", "1"},
			delay:     1,
			responses: responseDefault,
			expected:  "WARNING",
		},
		{
			name:      "critical",
			argv:      []string{"--critical", "1"},
			delay:     1,
			responses: responseDefault,
			expected:  "CRITICAL",
		},
		{
			name:      "timeout",
			argv:      []string{"--timeout", "-1"},
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
			go func() {
				if err := s.runServe(); err != nil {
					t.Fatal(err.Error())
				}
			}()
			port := strings.Split(ln.Addr().String(), ":")[1]
			argv := []string{"--host", "127.0.0.1", "--port", port}
			argv = append(argv, c.argv...)
			ckr := run(argv)
			if ckr.Status.String() != c.expected {
				t.Errorf("%s: %s", ckr.Status.String(), ckr.Message)
			}
		})
	}
}
