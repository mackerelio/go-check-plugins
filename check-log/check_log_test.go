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

	stateFile := getStateFile(opts.StateDir, logf)

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

	stateFile := getStateFile(opts.StateDir, logf)

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
