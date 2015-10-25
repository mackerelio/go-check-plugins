package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStateFile(t *testing.T) {
	sPath := getStateFile("/var/lib", "C:/Windows/hoge")
	assert.Equal(t, sPath, "/var/lib/C/Windows/hoge", "drive letter should be cared")

	sPath = getStateFile("/var/lib", "/linux/hoge")
	assert.Equal(t, sPath, "/var/lib/linux/hoge", "drive letter should be cared")
}

func TestWriteBytesToSkip(t *testing.T) {
	f := ".tmp/fuga/piyo"
	err := writeBytesToSkip(f, 15)
	assert.Equal(t, err, nil, "err should be nil")

	skipBytes, err := getBytesToSkip(f)
	assert.Equal(t, err, nil, "err should be nil")
	assert.Equal(t, skipBytes, int64(15))
}

func TestSearchReader(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Fatalf("TempDir failed: %s", err)
	}
	defer os.RemoveAll(dir)

	opts := &logOpts{
		StateDir: dir,
		LogFile:  filepath.Join(dir, "dummy"),
		Pattern:  `FATAL`,
	}
	opts.prepare()

	content := `FATAL 11
OK
FATAL 22
Fatal
`
	r := strings.NewReader(content)
	warnNum, critNum, readBytes, errLines, err := opts.searchReader(r)

	assert.Equal(t, int64(2), warnNum, "warnNum should be 2")
	assert.Equal(t, int64(2), critNum, "critNum should be 2")
	assert.Equal(t, "FATAL 11\nFATAL 22\n", errLines, "invalid errLines")
	assert.Equal(t, int64(len(content)), readBytes, "readBytes should be 26")
}

func TestSearchReaderWithLevel(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf := filepath.Join(dir, "dummy")
	ptn := `FATAL level:([0-9]+)`
	opts, _ := parseArgs([]string{"-s", dir, "-f", logf, "-i", "-p", ptn, "--critical-level=17", "--warning-level=11"})
	if !reflect.DeepEqual(&logOpts{
		StateDir:        dir,
		LogFile:         filepath.Join(dir, "dummy"),
		CaseInsensitive: true,
		Pattern:         `FATAL level:([0-9]+)`,
		WarnLevel:       11,
		CritLevel:       17,
	}, opts) {
		t.Errorf("something went wrong")
	}
	opts.prepare()

	content := `FATAL level:11
OK
FATAL level:22
Fatal level:17
`
	r := strings.NewReader(content)
	warnNum, critNum, readBytes, errLines, err := opts.searchReader(r)

	assert.Equal(t, int64(2), warnNum, "warnNum should be 2")
	assert.Equal(t, int64(1), critNum, "critNum should be 1")
	assert.Equal(t, "FATAL level:22\nFatal level:17\n", errLines, "invalid errLines")
	assert.Equal(t, int64(len(content)), readBytes, "readBytes should be 26")
}
