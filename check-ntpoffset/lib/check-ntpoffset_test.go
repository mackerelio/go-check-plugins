package checkntpoffset

import (
	"strings"
	"testing"
)

func TestParseNTPOffsetFromChrony(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect float64
	}{
		{
			name: "normal",
			input: `Reference ID    : 160.16.75.242 (sv01.azsx.net)
Stratum         : 3
Ref time (UTC)  : Thu May  4 11:51:30 2017
System time     : 0.000033190 seconds slow of NTP time
Last offset     : +0.000003614 seconds
RMS offset      : 0.000017540 seconds
Frequency       : 10.880 ppm fast
Residual freq   : -0.000 ppm
Skew            : 0.003 ppm
Root delay      : 0.003541 seconds
Root dispersion : 0.000849 seconds
Update interval : 1030.4 seconds
Leap status     : Normal
`,
			expect: 0.003614,
		},
	}

	for _, tc := range testCases {
		offset, err := parseNTPOffsetFromChrony(strings.NewReader(tc.input), false)
		if err != nil {
			t.Fatalf("error should be nil but got: %v", err)
		}
		if offset != tc.expect {
			t.Errorf("invalid offset: %f (expected: %f)", offset, tc.expect)
		}
	}
}

func TestParseNTPOffsetFromNTPD(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect float64
	}{
		{
			name:   "normal",
			input:  "stratum=3, offset=0.504\n",
			expect: 0.504,
		},
		{
			name: "ntp on 4.2.2p1-18.el5.centos",
			input: `assID=0 status=06f4 leap_none, sync_ntp, 15 events, event_peer/strat_chg,
stratum=3, offset=0.180
`,
			expect: 0.18,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			offset, err := parseNTPOffsetFromNTPD(strings.NewReader(tc.input), false)
			if err != nil {
				t.Fatalf("error should be nil but got: %v", err)
			}
			if offset != tc.expect {
				t.Errorf("%s: invalid offset: %f (expected: %f)", tc.name, offset, tc.expect)
			}
		})
	}
}
