package checkprocs

import (
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

func TestOptimizeStatus(t *testing.T) {
	var CritOver int64 = 100
	var WarningOver int64 = 80
	var WarningUnder int64 = 40
	var CritUnder int64 = 10

	opts.CritOver = &CritOver
	opts.WarningOver = &WarningOver
	opts.WarningUnder = WarningUnder
	opts.CritUnder = CritUnder

	assert.Equal(t, checkers.OK, optimizeStatus(80, checkers.OK))
	assert.Equal(t, checkers.WARNING, optimizeStatus(81, checkers.OK))
	assert.Equal(t, checkers.WARNING, optimizeStatus(100, checkers.OK))
	assert.Equal(t, checkers.CRITICAL, optimizeStatus(101, checkers.OK))
	assert.Equal(t, checkers.OK, optimizeStatus(40, checkers.OK))
	assert.Equal(t, checkers.WARNING, optimizeStatus(39, checkers.OK))
	assert.Equal(t, checkers.WARNING, optimizeStatus(10, checkers.OK))
	assert.Equal(t, checkers.CRITICAL, optimizeStatus(9, checkers.OK))
}
