package checkping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIPv6(t *testing.T) {
	testCases := []struct {
		casename          string
		host              string
		expectDetectation bool
	}{
		{
			casename:          "IPv4 IP address",
			host:              "127.0.0.1",
			expectDetectation: false,
		},
		{
			casename:          "IPv6 IP address 01",
			host:              "fe80::1",
			expectDetectation: true,
		},
		{
			casename:          "IPv6 IP address 02",
			host:              "2001:db8::1",
			expectDetectation: true,
		},
		{
			casename:          "IPv6 IP address 03",
			host:              "::1",
			expectDetectation: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.casename, func(t *testing.T) {
			result := isIPv6(tc.host)
			assert.Equal(t, tc.expectDetectation, result, "something went wrong")
		})
	}
}
