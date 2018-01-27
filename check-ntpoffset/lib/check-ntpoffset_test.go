package checkntpoffset

import (
	"strings"
	"testing"
)

func TestParseNTPOffsetFromChrony(t *testing.T) {
	offset, err := parseNTPOffsetFromChrony(strings.NewReader(
		`Reference ID    : 160.16.75.242 (sv01.azsx.net)
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
`))
	if err != nil {
		t.Fatalf("error should be nil but got: %v", err)
	}
	expect := 0.003614
	if offset != expect {
		t.Errorf("invalid offset: %f (expected: %f)", offset, expect)
	}
}
