package checkprocs

import (
	"runtime"
	"testing"
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
