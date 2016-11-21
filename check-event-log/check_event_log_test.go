package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/stretchr/testify/assert"
)

func TestGetStateFile(t *testing.T) {
	sPath := getStateFile("/var/lib", "C:/Windows/hoge")
	assert.Equal(t, sPath, "/var/lib/C/Windows/hoge", "drive letter should be cared")

	sPath = getStateFile("/var/lib", "/linux/hoge")
	assert.Equal(t, sPath, "/var/lib/linux/hoge", "drive letter should be cared")
}

func TestWriteLastOffset(t *testing.T) {
	f := ".tmp/fuga/piyo"
	err := writeLastOffset(f, 15)
	assert.Equal(t, err, nil, "err should be nil")

	recordNumber, err := getLastOffset(f)
	assert.Equal(t, err, nil, "err should be nil")
	assert.Equal(t, recordNumber, int64(15))
}

func raiseEvent(t *testing.T, typ int, msg string) {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unk, err := oleutil.CreateObject("Wscript.Shell")
	if err != nil {
		t.Fatal(err)
	}
	disp, err := unk.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		t.Fatal(err)
	}
	oleutil.MustCallMethod(disp, "LogEvent", typ, msg)
}

func TestRun(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-event-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	opts, _ := parseArgs([]string{"-s", dir, "--log", "Application"})
	opts.prepare()

	stateFile := getStateFile(opts.StateDir, "Application")

	recordNumber, _ := getLastOffset(stateFile)
	lastNumber := recordNumber
	assert.Equal(t, int64(0), recordNumber, "something went wrong")

	testEmpty := func() {
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, 0, recordNumber, "something went wrong")
	}
	testEmpty()

	lastNumber = recordNumber

	testInfo := func() {
		raiseEvent(t, 0, "check-event-log: something info occured")
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
	}
	testInfo()

	lastNumber = recordNumber

	testWarning := func() {
		raiseEvent(t, 2, "check-event-log: something warning occured")
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
	}
	testWarning()

	lastNumber = recordNumber

	testError := func() {
		raiseEvent(t, 1, "check-event-log: something error occured")
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(1), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
	}
	testError()

	lastNumber = recordNumber

	opts, _ = parseArgs([]string{"-s", dir, "--log", "Application", "-r"})
	opts.prepare()

	testReturn := func() {
		raiseEvent(t, 1, "check-event-log: something error occured")
		raiseEvent(t, 2, "check-event-log: something warning occured")
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "something went wrong")
		assert.Equal(t, int64(1), c, "something went wrong")
		assert.Equal(t, "WSH:check-event-log: something error occured\nWSH:check-event-log: something warning occured\n", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
	}
	testReturn()
}

func TestSourcePattern(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-event-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	opts, _ := parseArgs([]string{"-s", dir, "--log", "Application"})
	opts.prepare()

	stateFile := getStateFile(opts.StateDir, "Application")

	recordNumber, _ := getLastOffset(stateFile)
	lastNumber := recordNumber
	assert.Equal(t, int64(0), recordNumber, "something went wrong")

	testEmpty := func() {
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, 0, recordNumber, "something went wrong")
	}
	testEmpty()

	lastNumber = recordNumber

	opts, _ = parseArgs([]string{"-s", dir, "--log", "Application", "--message-pattern", "テストエラーが(発生しました|起きました)"})
	opts.prepare()

	testMessagePattern := func() {
		raiseEvent(t, 1, "check-event-log: テストエラーが発生しました")
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(1), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
	}
	testMessagePattern()

	opts, _ = parseArgs([]string{"-s", dir, "--log", "Application", "--source-pattern", "[Ww][Ss][Hh]"})
	opts.prepare()

	testSourcePattern := func() {
		raiseEvent(t, 2, "check-event-log: テストエラーが発生しました")
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
	}
	testSourcePattern()
}

func TestFailFirst(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-event-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	opts, _ := parseArgs([]string{"-s", dir, "--log", "Application", "--fail-first", "--warning-over", "0", "--critical-over", "0"})
	opts.prepare()

	testFailFirst := func() {
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.NotEqual(t, int64(0), w, "something went wrong")
		assert.NotEqual(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")
	}
	testFailFirst()
}
