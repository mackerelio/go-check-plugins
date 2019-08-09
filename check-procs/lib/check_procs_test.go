package checkprocs

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestProcs(t *testing.T) {
	procs, err := getProcs()
	if err != nil {
		t.Fatal(err)
	}
	for _, proc := range procs {
		if proc.cmd != "go" {
			continue
		}
		if proc.user == "" && runtime.GOOS != "windows" {
			t.Fatal("user should not be empty")
		}
		if proc.ppid == "" && runtime.GOOS != "windows" {
			t.Fatal("ppid should not be empty")
		}
		if proc.pid == "" {
			t.Fatal("pid should not be empty")
		}
		if proc.vsz == 0 {
			t.Fatal("vsz should not be 0")
		}
		if proc.rss == 0 {
			t.Fatal("rss should not be 0")
		}
		if proc.thcount == 0 {
			t.Fatal("thcount should not be 0")
		}
		if proc.state == "" && runtime.GOOS != "windows" {
			t.Fatal("state should not be empty")
		}
		if proc.esec == 0 {
			t.Fatal("esec should not be 0")
		}
		if proc.csec == 0 && runtime.GOOS != "windows" {
			t.Fatal("csec should not be 0")
		}
	}
}

func TestMergeStatus(t *testing.T) {
	var CritOver int64 = 100
	var WarningOver int64 = 80
	var WarningUnder int64 = 40
	var CritUnder int64 = 10

	opts.CritOver = &CritOver
	opts.WarningOver = &WarningOver
	opts.WarningUnder = WarningUnder
	opts.CritUnder = CritUnder

	assert.Equal(t, checkers.OK, mergeStatus(80, checkers.OK))
	assert.Equal(t, checkers.WARNING, mergeStatus(81, checkers.OK))
	assert.Equal(t, checkers.WARNING, mergeStatus(100, checkers.OK))
	assert.Equal(t, checkers.CRITICAL, mergeStatus(101, checkers.OK))
	assert.Equal(t, checkers.OK, mergeStatus(40, checkers.OK))
	assert.Equal(t, checkers.WARNING, mergeStatus(39, checkers.OK))
	assert.Equal(t, checkers.WARNING, mergeStatus(10, checkers.OK))
	assert.Equal(t, checkers.CRITICAL, mergeStatus(9, checkers.OK))
}

func TestGatherMsg(t *testing.T) {
	var count int64 = 1
	opts.CmdPatterns = []string{"foo", "bar"}

	for _, pattern := range opts.CmdPatterns {
		expected := fmt.Sprintf("Found %d matching processes; cmd /%s/", count, pattern)
		assert.Equal(t, expected, gatherMsg(count, pattern))
	}
}
