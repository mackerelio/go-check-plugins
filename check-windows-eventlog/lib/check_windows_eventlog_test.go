// +build windows

package checkwindowseventlog

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/stretchr/testify/assert"
)

func TestGetStateFile(t *testing.T) {
	opts := &logOpts{
		StateDir: "/var/lib",
		origArgs: []string{},
	}
	opts.prepare()
	sPath := opts.getStateFile("Application")
	if runtime.GOOS == "windows" {
		sPath = filepath.ToSlash(sPath)
	}
	assert.Equal(t, sPath, "/var/lib/Application-d41d8cd98f00b204e9800998ecf8427e", "arguments should be cared")

	opts = &logOpts{
		StateDir: "/var/lib",
		origArgs: []string{"foo", "bar"},
	}
	opts.prepare()
	sPath = opts.getStateFile("Security")
	if runtime.GOOS == "windows" {
		sPath = filepath.ToSlash(sPath)
	}
	assert.Equal(t, sPath, "/var/lib/Security-327b6f07435811239bc47e1544353273", "arguments should be cared")
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
	dir, err := ioutil.TempDir("", "check-windows-eventlog-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	origArgs := []string{"-s", dir, "--log", "Application"}
	args := make([]string, len(origArgs))
	copy(args, origArgs)
	opts, err := parseArgs(args)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(&logOpts{
		StateDir: dir,
		Log:      `Application`,
		origArgs: origArgs,
	}, opts) {
		t.Errorf("something went wrong: %#v", opts)
	}

	opts.prepare()

	stateFile := opts.getStateFile("Application")

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
		raiseEvent(t, 0, "check-windows-eventlog: something info occured")
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
		raiseEvent(t, 2, "check-windows-eventlog: something warning occured")
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
		raiseEvent(t, 1, "check-windows-eventlog: something error occured")
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

	origArgs = []string{"-s", dir, "--log", "Application", "-r"}
	args = make([]string, len(origArgs))
	copy(args, origArgs)
	opts, err = parseArgs(args)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(&logOpts{
		StateDir:      dir,
		Log:           `Application`,
		ReturnContent: true,
		origArgs:      origArgs,
	}, opts) {
		t.Errorf("something went wrong: %#v", opts)
	}

	opts.prepare()

	stateFile = opts.getStateFile("Application")

	recordNumber, _ = getLastOffset(stateFile)
	lastNumber = recordNumber
	assert.Equal(t, int64(0), recordNumber, "something went wrong")

	testEmpty = func() {
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, 0, recordNumber, "something went wrong")
	}
	testEmpty()

	testReturn := func() {
		raiseEvent(t, 1, "check-windows-eventlog: something error occured")
		raiseEvent(t, 2, "check-windows-eventlog: something warning occured")
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "something went wrong")
		assert.Equal(t, int64(1), c, "something went wrong")
		assert.Equal(t, "WSH:check-windows-eventlog: something error occured\nWSH:check-windows-eventlog: something warning occured\n", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
	}
	testReturn()
}

func TestPattern(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-windows-eventlog-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	opts, _ := parseArgs([]string{"-s", dir, "--log", "Application"})
	opts.prepare()

	stateFile := opts.getStateFile("Application")
	recordNumber, _ := getLastOffset(stateFile)
	lastNumber := recordNumber
	assert.Equal(t, int64(0), recordNumber, "something went wrong")

	reset := func(args []string) {
		opts, _ = parseArgs(args)
		opts.prepare()

		lastNumber = recordNumber
		stateFile = opts.getStateFile("Application")
		writeLastOffset(stateFile, lastNumber)
	}

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

	reset([]string{"-s", dir, "--log", "Application", "--message-pattern", "テストエラーが(発生しました|起きました)"})

	testMessagePattern := func() {
		raiseEvent(t, 1, "check-windows-eventlog: テストエラーが発生しました")
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(1), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
	}
	testMessagePattern()

	reset([]string{"-s", dir, "--log", "Application", "--source-pattern", "[Ww][Ss][Hh]"})

	testSourcePattern := func() {
		raiseEvent(t, 2, "check-windows-eventlog: テストエラーが発生しました")
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
	}
	testSourcePattern()

	reset([]string{"-s", dir, "--log", "Application", "--message-pattern", "テストエラー", "--message-exclude", "起きました"})

	testMessageExclude := func() {
		raiseEvent(t, 1, "check-windows-eventlog: テストエラーが発生しました")
		raiseEvent(t, 1, "check-windows-eventlog: テストエラーが起きました")
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(1), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
	}
	testMessageExclude()

	reset([]string{"-s", dir, "--log", "Application", "--source-pattern", "[Ww][Ss]", "--source-exclude", "[h]"})

	testSourceExclude := func() {
		raiseEvent(t, 2, "check-windows-eventlog: テストエラーが発生しました")
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
	}
	testSourceExclude()

	reset([]string{"-s", dir, "--log", "Application", "--source-pattern", "[Ww][Ss]", "--source-exclude", "[H]"})

	testSourceExclude = func() {
		raiseEvent(t, 2, "check-windows-eventlog: テストエラーが発生しました")
		w, c, errLines, err := opts.searchLog("Application")
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		recordNumber, _ = getLastOffset(stateFile)
		assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
	}
	testSourceExclude()
}

func TestIDs(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-windows-eventlog-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	opts, _ := parseArgs([]string{"-s", dir, "--log", "Application"})
	opts.prepare()

	stateFile := opts.getStateFile("Application")
	recordNumber, _ := getLastOffset(stateFile)
	lastNumber := recordNumber
	assert.Equal(t, int64(0), recordNumber, "something went wrong")

	reset := func(args []string) {
		opts, _ = parseArgs(args)
		opts.prepare()

		lastNumber = recordNumber
		stateFile = opts.getStateFile("Application")
		writeLastOffset(stateFile, lastNumber)
	}

	_, _, _, err = opts.searchLog("Application")
	if err != nil {
		t.Fatal(err)
	}
	recordNumber, _ = getLastOffset(stateFile)

	tests := []struct {
		pattern string
		exclude string
		c       int64
		w       int64
	}{
		{"1,2", "", 1, 1},
		{"1", "1", 0, 0},
		{"1", "2", 1, 0},
		{"0,1-2", "", 1, 1},
		{"1-2", "2", 1, 0},
		{"1", "2", 1, 0},
		{"0", "2", 0, 0},
		{"", "1-2", 0, 0},
		{"", "2-3", 1, 0},
		{"", "3-2", 1, 0},
		{"2", "1-2", 0, 0},
		{"2", "2-3", 0, 0},
		{"2", "3-4", 0, 1},
		{"2", "1,3", 0, 1},
	}

	for _, test := range tests {
		reset([]string{
			"-s", dir, "--log", "Application", "--source-pattern", "WSH",
			"--event-id-pattern", test.pattern,
			"--event-id-exclude", test.exclude,
		})

		testID := func() {
			t.Logf("testing: --event-id-pattern %q, --event-id-exclude=%q", test.pattern, test.exclude)
			raiseEvent(t, 1, "check-windows-eventlog: critical event occured")
			raiseEvent(t, 2, "check-windows-eventlog: warning event occured")
			w, c, errLines, err := opts.searchLog("Application")
			assert.Equal(t, err, nil, "err should be nil")
			assert.Equal(t, int64(test.w), w, "something went wrong")
			assert.Equal(t, int64(test.c), c, "something went wrong")
			assert.Equal(t, "", errLines, "something went wrong")

			recordNumber, _ = getLastOffset(stateFile)
			assert.NotEqual(t, lastNumber, recordNumber, "something went wrong")
		}
		testID()
	}
}

func TestFailFirst(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-windows-eventlog-test")
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
