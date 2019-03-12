package checkntpoffset

import (
	"strings"
	"testing"
)

func TestParseNTPOffsetFromChrony(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		checkStratum bool
		expect       float64
		expectError  string
	}{
		{
			name: "normal (checkStratum = false)",
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
			checkStratum: false,
			expect:       0.003614,
		},
		{
			name: "normal (checkStratum = true, synchronized)",
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
			checkStratum: true,
			expect:       0.003614,
		},
		{
			name: "normal (checkStratum = true, unsynchronized)",
			input: `Reference ID    : 00000000 ()
Stratum         : 0
Ref time (UTC)  : Thu Jan 01 00:00:00 1970
System time     : 0.000000000 seconds fast of NTP time
Last offset     : +0.000000000 seconds
RMS offset      : 0.000000000 seconds
Frequency       : 281.118 ppm slow
Residual freq   : +0.000 ppm
Skew            : 0.000 ppm
Root delay      : 1.000000 seconds
Root dispersion : 1.000000 seconds
Update interval : 0.0 seconds
Leap status     : Not synchronised
`,
			checkStratum: true,
			expectError:  "not synchronized to stratum",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			offset, err := parseNTPOffsetFromChrony(strings.NewReader(tc.input), tc.checkStratum)
			if tc.expectError != "" {
				if err == nil {
					t.Error("error should not be nil")
				}
				if err.Error() != tc.expectError {
					t.Errorf("unexpected error: %s (expected: %s)", err.Error(), tc.expectError)
				}
			} else {
				if err != nil {
					t.Fatalf("error should be nil but got: %v", err)
				}
				if offset != tc.expect {
					t.Errorf("invalid offset: %f (expected: %f)", offset, tc.expect)
				}
			}
		})
	}
}

func TestParseNTPOffsetFromNTPD(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		checkStratum bool
		expect       float64
		expectError  string
	}{
		{
			name:         "normal (checkStratum = false)",
			input:        "stratum=3, offset=0.504\n",
			checkStratum: false,
			expect:       0.504,
		},
		{
			name:         "normal (checkStratum = true, synchronized)",
			input:        "stratum=3, offset=0.504\n",
			checkStratum: true,
			expect:       0.504,
		},
		{
			name:         "normal (checkStratum = true, unsynchronized)",
			input:        "stratum=16, offset=0.000000\n",
			checkStratum: true,
			expectError:  "not synchronized to stratum",
		},
		{
			name: "ntp on 4.2.2p1-18.el5.centos (checkStratum = false)",
			input: `assID=0 status=06f4 leap_none, sync_ntp, 15 events, event_peer/strat_chg,
stratum=3, offset=0.180
`,
			checkStratum: false,
			expect:       0.18,
		},
		{
			name: "ntp on 4.2.2p1-18.el5.centos (checkStratum = true, synchronized)",
			input: `assID=0 status=06f4 leap_none, sync_ntp, 15 events, event_peer/strat_chg,
stratum=3, offset=0.180
`,
			checkStratum: true,
			expect:       0.18,
		},
		{
			name: "ntp on 4.2.2p1-18.el5.centos (checkStratum = true, unsynchronized)",
			input: `assID=0 status=06f4 leap_none, sync_ntp, 15 events, event_peer/strat_chg,
stratum=16, offset=0.000
`,
			checkStratum: true,
			expectError:  "not synchronized to stratum",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			offset, err := parseNTPOffsetFromNTPD(strings.NewReader(tc.input), tc.checkStratum)
			if tc.expectError != "" {
				if err == nil {
					t.Error("error should not be nil")
				}
				if err.Error() != tc.expectError {
					t.Errorf("unexpected error: %s (expected: %s)", err.Error(), tc.expectError)
				}
			} else {
				if err != nil {
					t.Fatalf("error should be nil but got: %v", err)
				}
				if offset != tc.expect {
					t.Errorf("invalid offset: %f (expected: %f)", offset, tc.expect)
				}
			}
		})
	}
}
