package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestGetStateFile(t *testing.T) {
	sPath := getStateFile("/var/lib", "C:/Windows/hoge", []string{})
	assert.Equal(t, "/var/lib/C/Windows/hoge-d41d8cd98f00b204e9800998ecf8427e", sPath, "drive letter should be cared")

	sPath = getStateFile("/var/lib", "/linux/hoge", []string{})
	assert.Equal(t, "/var/lib/linux/hoge-d41d8cd98f00b204e9800998ecf8427e", sPath, "arguments should be cared")

	sPath = getStateFile("/var/lib", "/linux/hoge", []string{"aa", "BB"})
	assert.Equal(t, "/var/lib/linux/hoge-c508092e97c59149a8644827e066ca83", sPath, "arguments should be cared")
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

func TestRun(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	defer fh.Close()

	ptn := `FATAL`
	opts, _ := parseArgs([]string{"-s", dir, "-f", logf, "-p", ptn})
	opts.prepare()

	stateFile := getStateFile(opts.StateDir, logf, opts.origArgs)

	bytes, _ := getBytesToSkip(stateFile)
	assert.Equal(t, int64(0), bytes, "something went wrong")

	testEmpty := func() {
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(0), bytes, "something went wrong")
	}
	testEmpty()

	l1 := "FATAL\nFATAL\n"
	test2Line := func() {
		fh.WriteString(l1)
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(2), w, "something went wrong")
		assert.Equal(t, int64(2), c, "something went wrong")
		assert.Equal(t, l1, errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(l1)), bytes, "something went wrong")
	}
	test2Line()

	testReadAgain := func() {
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(l1)), bytes, "something went wrong")
	}
	testReadAgain()

	l2 := "SUCCESS\n"
	testRecover := func() {
		fh.WriteString(l2)
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(l1)+len(l2)), bytes, "something went wrong")
	}
	testRecover()

	testSuccessAgain := func() {
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(l1)+len(l2)), bytes, "something went wrong")
	}
	testSuccessAgain()

	testErrorAgain := func() {
		fh.WriteString(l1)
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(2), w, "something went wrong")
		assert.Equal(t, int64(2), c, "something went wrong")
		assert.Equal(t, l1, errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(l1)*2+len(l2)), bytes, "something went wrong")
	}
	testErrorAgain()

	testRecoverAgain := func() {
		fh.WriteString(l2)
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(l1)*2+len(l2)*2), bytes, "something went wrong")
	}
	testRecoverAgain()

	testRotate := func() {
		fh.Close()
		os.Remove(logf)
		fh, _ = os.Create(logf)

		fh.WriteString(l2)
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(l2)), bytes, "something went wrong")
	}
	testRotate()
}

func TestRunWithGlob(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf1 := filepath.Join(dir, "dummy1")
	fh1, _ := os.Create(logf1)
	defer fh1.Close()

	logf2 := filepath.Join(dir, "dummy2")
	fh2, _ := os.Create(logf2)
	defer fh2.Close()

	ptn := `FATAL`
	params := []string{dir, "-g", filepath.Join(dir, "dummy*"), "-p", ptn}
	opts, _ := parseArgs(params)
	opts.prepare()

	testSuccess := func() {
		ckr := run(params)
		assert.Equal(t, checkers.OK, ckr.Status, "ckr.Status should be OK")
	}
	testSuccess()

	errorLine := "FATAL\n"
	testCriticalOnce := func() {
		fh1.WriteString(errorLine)
		ckr := run(params)
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "ckr.Status should be CRITICAL")
	}
	testCriticalOnce()

	testRecover := func() {
		ckr := run(params)
		assert.Equal(t, checkers.OK, ckr.Status, "ckr.Status should be OK")
	}
	testRecover()

	testCriticalAgain := func() {
		fh2.WriteString(errorLine)
		ckr := run(params)
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "ckr.Status should be CRITICAL")
	}
	testCriticalAgain()

}

func TestRunWithMiddleOfLine(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	defer fh.Close()

	ptn := `FATAL`
	opts, _ := parseArgs([]string{"-s", dir, "-f", logf, "-p", ptn})
	opts.prepare()

	stateFile := getStateFile(opts.StateDir, logf, opts.origArgs)

	bytes, _ := getBytesToSkip(stateFile)
	assert.Equal(t, int64(0), bytes, "something went wrong")

	testMiddleOfLine := func() {
		fh.WriteString("FATA")
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(0), bytes, "something went wrong")
	}
	testMiddleOfLine()

	testFail := func() {
		fh.WriteString("L\nSUCC")
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "something went wrong")
		assert.Equal(t, int64(1), c, "something went wrong")
		assert.Equal(t, "FATAL\n", errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len("FATAL\n")), bytes, "something went wrong")
	}
	testFail()
}

func TestRunWithNoState(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	defer fh.Close()

	ptn := `FATAL`
	opts, _ := parseArgs([]string{"-s", dir, "-f", logf, "-p", ptn, "--no-state"})
	opts.prepare()

	fatal := "FATAL\n"
	test2Line := func() {
		fh.WriteString(fatal)
		fh.WriteString(fatal)
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(2), w, "something went wrong")
		assert.Equal(t, int64(2), c, "something went wrong")
		assert.Equal(t, strings.Repeat(fatal, 2), errLines, "something went wrong")
	}
	test2Line()

	test1LineAgain := func() {
		fh.WriteString(fatal)
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(3), w, "something went wrong")
		assert.Equal(t, int64(3), c, "something went wrong")
		assert.Equal(t, strings.Repeat(fatal, 3), errLines, "something went wrong")
	}
	test1LineAgain()
}

func TestSearchReaderWithLevel(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf := filepath.Join(dir, "dummy")
	ptn := `FATAL level:([0-9]+)`
	origArgs := []string{"-s", dir, "-f", logf, "-i", "-p", ptn, "--critical-level=17", "--warning-level=11"}
	args := make([]string, len(origArgs))
	copy(args, origArgs)
	opts, _ := parseArgs(args)
	if !reflect.DeepEqual(&logOpts{
		StateDir:        dir,
		LogFile:         filepath.Join(dir, "dummy"),
		CaseInsensitive: true,
		Pattern:         `FATAL level:([0-9]+)`,
		WarnLevel:       11,
		CritLevel:       17,
		Missing:         "UNKNOWN",
		origArgs:        origArgs,
	}, opts) {
		t.Errorf("something went wrong: %#v", opts)
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

func TestRunWithEncoding(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	defer fh.Close()

	opts, _ := parseArgs([]string{"-s", dir, "-f", logf, "-p", `エラー`, "--encoding", "euc-jp"})
	opts.prepare()

	fatal := "\xa5\xa8\xa5\xe9\xa1\xbc\n" // エラー
	testEncoding := func() {
		fh.Write([]byte(fatal))
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "something went wrong")
		assert.Equal(t, int64(1), c, "something went wrong")
		assert.Equal(t, "エラー\n", errLines, "something went wrong")
	}
	testEncoding()
}

func TestRunWithoutEncoding(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	defer fh.Close()

	opts, _ := parseArgs([]string{"-s", dir, "-f", logf, "-p", `エラー`})
	opts.prepare()

	fatal := "\xa5\xa8\xa5\xe9\xa1\xbc\nエラー\n" // エラー
	testWithoutEncoding := func() {
		fh.Write([]byte(fatal))
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "something went wrong")
		assert.Equal(t, int64(1), c, "something went wrong")
		assert.Equal(t, "エラー\n", errLines, "something went wrong")
	}
	testWithoutEncoding()

	fatal = "エラー\n"
	testWithEncoding := func() {
		fh.Write([]byte(fatal))
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "something went wrong")
		assert.Equal(t, int64(1), c, "something went wrong")
		assert.Equal(t, "エラー\n", errLines, "something went wrong")
	}
	testWithEncoding()
}

func TestRunWithMissingOk(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf := filepath.Join(dir, "dummy")

	ptn := `FATAL`
	missing := `OK`
	params := []string{"-s", dir, "-f", logf, "-p", ptn, "--missing", missing}
	opts, _ := parseArgs(params)
	opts.prepare()

	testRunLogFileMissing := func() {
		ckr := run(params)
		assert.Equal(t, ckr.Status, checkers.OK, "ckr.Status should be OK")
		msg := fmt.Sprintf("0 warnings, 0 criticals for pattern /FATAL/.\nThe following 1 files are missing.\n%s", logf)
		assert.Equal(t, ckr.Message, msg, "something went wrong")
	}
	testRunLogFileMissing()
}

func TestRunWithMissingWarning(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf := filepath.Join(dir, "dummy")

	ptn := `FATAL`
	missing := `WARNING`
	params := []string{"-s", dir, "-f", logf, "-p", ptn, "--missing", missing}
	opts, _ := parseArgs(params)
	opts.prepare()

	testRunLogFileMissing := func() {
		ckr := run(params)
		assert.Equal(t, ckr.Status, checkers.WARNING, "ckr.Status should be WARNING")
		msg := fmt.Sprintf("0 warnings, 0 criticals for pattern /FATAL/.\nThe following 1 files are missing.\n%s", logf)
		assert.Equal(t, ckr.Message, msg, "something went wrong")
	}
	testRunLogFileMissing()
}

func TestRunWithMissingCritical(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf := filepath.Join(dir, "dummy")

	ptn := `FATAL`
	missing := `CRITICAL`
	params := []string{"-s", dir, "-f", logf, "-p", ptn, "--missing", missing}
	opts, _ := parseArgs(params)
	opts.prepare()

	testRunLogFileMissing := func() {
		ckr := run(params)
		assert.Equal(t, ckr.Status, checkers.CRITICAL, "ckr.Status should be CRITICAL")
		msg := fmt.Sprintf("0 warnings, 0 criticals for pattern /FATAL/.\nThe following 1 files are missing.\n%s", logf)
		assert.Equal(t, ckr.Message, msg, "something went wrong")
	}
	testRunLogFileMissing()
}

func TestRunWithMissingUnknown(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf := filepath.Join(dir, "dummy")

	ptn := `FATAL`
	missing := `UNKNOWN`
	params := []string{"-s", dir, "-f", logf, "-p", ptn, "--missing", missing}
	opts, _ := parseArgs(params)
	opts.prepare()

	testRunLogFileMissing := func() {
		ckr := run(params)
		assert.Equal(t, ckr.Status, checkers.UNKNOWN, "ckr.Status should be UNKNOWN")
		msg := fmt.Sprintf("0 warnings, 0 criticals for pattern /FATAL/.\nThe following 1 files are missing.\n%s", logf)
		assert.Equal(t, ckr.Message, msg, "something went wrong")
	}
	testRunLogFileMissing()
}
