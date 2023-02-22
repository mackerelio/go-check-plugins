package checklog

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestGetStateFile(t *testing.T) {
	sPath := getStateFile("/var/lib", "C:/Windows/hoge", []string{})
	assert.Equal(t, "/var/lib/C/Windows/hoge-d41d8cd98f00b204e9800998ecf8427e.json", filepath.ToSlash(sPath), "drive letter should be cared")

	sPath = getStateFile("/var/lib", "/linux/hoge", []string{})
	assert.Equal(t, "/var/lib/linux/hoge-d41d8cd98f00b204e9800998ecf8427e.json", filepath.ToSlash(sPath), "arguments should be cared")

	sPath = getStateFile("/var/lib", "/linux/hoge", []string{"aa", "BB"})
	assert.Equal(t, "/var/lib/linux/hoge-c508092e97c59149a8644827e066ca83.json", filepath.ToSlash(sPath), "arguments should be cared")
}

func TestSaveState(t *testing.T) {
	f := ".tmp/fuga/piyo.json"
	err := saveState(f, &state{SkipBytes: 15, Inode: 150})
	t.Cleanup(func() {
		os.RemoveAll(".tmp")
	})
	assert.Equal(t, err, nil, "err should be nil")

	state, err := loadState(f)
	assert.Equal(t, err, nil, "err should be nil")
	assert.Equal(t, state.SkipBytes, int64(15))
	assert.Equal(t, state.Inode, uint(150))
}

func TestGetInode(t *testing.T) {
	f := ".tmp/hoge/piyo.json"
	state := &state{SkipBytes: 15, Inode: 150}
	t.Cleanup(func() {
		os.RemoveAll(".tmp")
	})

	i, err := getInode(f)
	assert.Equal(t, err, nil, "err should be nil")
	assert.Equal(t, i, uint(0), "inode should be empty")

	saveState(f, state)

	i, err = getInode(f)
	assert.Equal(t, err, nil, "err should be nil")
	assert.Equal(t, state.Inode, uint(150))
}

func TestGetBytesToSkip(t *testing.T) {
	// fallback testing for backward compatibility
	oldf := ".tmp/fuga/piyo"
	newf := ".tmp/fuga/piyo.json"
	state := &state{SkipBytes: 15}
	os.MkdirAll(filepath.Dir(oldf), 0755)
	ioutil.WriteFile(oldf, []byte(fmt.Sprintf("%d", state.SkipBytes)), 0600)
	t.Cleanup(func() {
		os.RemoveAll(".tmp")
	})

	n, err := getBytesToSkip(newf)
	assert.Equal(t, err, nil, "err should be nil")
	assert.Equal(t, state.SkipBytes, n)

	saveState(newf, state)

	n, err = getBytesToSkip(newf)
	assert.Equal(t, err, nil, "err should be nil")
	assert.Equal(t, state.SkipBytes, n)
}

func TestSearchReader(t *testing.T) {
	dir := t.TempDir()

	opts := &logOpts{
		StateDir: dir,
		LogFile:  filepath.Join(dir, "dummy"),
		Pattern:  []string{`FATAL`},
	}
	opts.prepare()

	content := `FATAL 11
OK
FATAL 22
Fatal
`
	r := strings.NewReader(content)
	warnNum, critNum, readBytes, errLines, _ := opts.searchReader(context.Background(), r)

	assert.Equal(t, int64(2), warnNum, "warnNum should be 2")
	assert.Equal(t, int64(2), critNum, "critNum should be 2")
	assert.Equal(t, "FATAL 11\nFATAL 22\n", errLines, "invalid errLines")
	assert.Equal(t, int64(len(content)), readBytes, "readBytes should be 26")
}

func TestScenario(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	t.Cleanup(func() {
		fh.Close()
	})

	ptn := `FATAL`
	opts, _ := parseArgs([]string{"-s", dir, "-f", logf, "-p", ptn})
	opts.prepare()

	stateFile := getStateFile(opts.StateDir, logf, opts.origArgs)

	bytes, _ := getBytesToSkip(stateFile)
	assert.Equal(t, int64(0), bytes, "should be a 0-byte indicated value")

	t.Run("read an empty file", func(t *testing.T) {
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "should not be detected as it is empty")
		assert.Equal(t, int64(0), c, "should not be detected as it is empty")
		assert.Equal(t, "", errLines, "should not be detected as it is empty")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(0), bytes, "should be a 0-byte indicated value, because file is empty")
	})

	linesOf2Fatals := "FATAL\nFATAL\n"
	t.Run("give a detection string and generate an error condition", func(t *testing.T) {
		fh.WriteString(linesOf2Fatals)

		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(2), w, "there are two error strings, so two should be detected")
		assert.Equal(t, int64(2), c, "there are two error strings, so two should be detected")
		assert.Equal(t, linesOf2Fatals, errLines, "given string should be detected")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(linesOf2Fatals)), bytes, "the pointer should be moved by the amount of the character string written in the file")
	})

	linesOf3Fatals := "FATAL\nFATAL\nFATAL\n"
	t.Run("change the state of the detected string and check that the error state changes", func(t *testing.T) {
		fh.WriteString(linesOf3Fatals)
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(3), w, "there are three error strings, so three should be detected")
		assert.Equal(t, int64(3), c, "there are three error strings, so three should be detected")
		assert.Equal(t, linesOf3Fatals, errLines, "given string should be detected")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(linesOf2Fatals)+len(linesOf3Fatals)), bytes, "the pointer should be moved by the amount of the character string written in the file")
	})

	lineOfSuccess := "SUCCESS\n"
	t.Run("change the state of the detected string and check that the error state is resolved.", func(t *testing.T) {
		fh.WriteString(lineOfSuccess)
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "should not be detected as the error condition has been resolved")
		assert.Equal(t, int64(0), c, "should not be detected as the error condition has been resolved")
		assert.Equal(t, "", errLines, "should not be detected as the error condition has been resolved")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(linesOf2Fatals)+len(linesOf3Fatals)+len(lineOfSuccess)), bytes, "the pointer should be moved by the amount of the character string written in the file")
	})

	t.Run("when the state of the file does not change, the state continues", func(t *testing.T) {
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "should not be detected as there is no change in state")
		assert.Equal(t, int64(0), c, "should not be detected as there is no change in state")
		assert.Equal(t, "", errLines, "should not be detected as there is no change in state")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(linesOf2Fatals)+len(linesOf3Fatals)+len(lineOfSuccess)), bytes, "the pointer should be moved by the amount of the character string written in the file")
	})

	t.Run("detected when there is a change in the state of the file, resulting in an error condition", func(t *testing.T) {
		fh.WriteString(linesOf3Fatals)
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(3), w, "there are three error strings, so three should be detected")
		assert.Equal(t, int64(3), c, "there are three error strings, so three should be detected")
		assert.Equal(t, linesOf3Fatals, errLines, "given string should be detected")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(linesOf2Fatals)+len(linesOf3Fatals)*2+len(lineOfSuccess)), bytes, "the pointer should be moved by the amount of the character string written in the file")
	})

	t.Run("detected when there is a change in the state of the file, the error has been resolved", func(t *testing.T) {
		fh.WriteString(lineOfSuccess)
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "should not be detected as the error condition has been resolved")
		assert.Equal(t, int64(0), c, "should not be detected as the error condition has been resolved")
		assert.Equal(t, "", errLines, "should not be detected as the error condition has been resolved")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(linesOf2Fatals)+len(linesOf3Fatals)*2+len(lineOfSuccess)*2), bytes, "the pointer should be moved by the amount of the character string written in the file")
	})

	t.Run("detect that a file has been rotated", func(t *testing.T) {
		// delete the inherited file and create a new file with the same name but a different inode.
		fh.Close()
		os.Remove(logf)
		fh, _ = os.Create(logf)

		fh.WriteString(linesOf3Fatals)
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(3), w, "the file is being rotated, so an error condition should be detected")
		assert.Equal(t, int64(3), c, "the file is being rotated, so an error condition should be detected")
		assert.Equal(t, linesOf3Fatals, errLines, "the file is being rotated, so an error condition should be detected")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(linesOf3Fatals)), bytes, "the file is being rotated, the pointer should be the size of the new file")
	})

	// Should test that check-log stops reading logs when timed out.
	// If a period (10*time.Millisecond in below) is very short,
	// normal behavior such as open a file, read it, etc could reach over a period.
	opts.testHookNewBufferedReader = func(r io.Reader) *bufio.Reader {
		return bufio.NewReaderSize(&slowReader{
			r: r,
			d: 10 * time.Millisecond,
			n: 1,
		}, 1)
	}
	t.Run("cancel", func(t *testing.T) {
		// This test checks searchLog keeps reading until the first EOL even if ctx is cancelled.
		// To guarantee to read at least once, a timeout sec
		// should be choice it is greater than reading the state file.
		fh.WriteString("FATAL\nFATAL\nFATAL\n")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
		t.Cleanup(func() {
			cancel()
		})

		w, c, errLines, err := opts.searchLog(ctx, logf)

		// searchLog should read only first line, so the result counts only the first `FATAL`
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "there are one error strings, so one should be detected")
		assert.Equal(t, int64(1), c, "there are one error strings, so one should be detected")
		assert.Equal(t, "FATAL\n", errLines, "the first `FATAL` should be detected")
	})
	opts.testHookNewBufferedReader = nil

	t.Run("cancel before processing", func(t *testing.T) {
		switch runtime.GOOS {
		case "windows":
			// TODO(lufia): Is there a file that a user running `go test` can't read on Windows?
			t.Skip()
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		cmdline := []string{"--no-state", "-p", "FATAL", "-f", "/etc/sudoers"}
		result := run(ctx, cmdline)
		assert.Equal(t, checkers.OK, result.Status, "OK should be detected")
		assert.Equal(t, "0 warnings, 0 criticals for pattern /FATAL/.", result.Message, "message with content where `FATAL` is not detected")
	})

	opts.testHookNewBufferedReader = func(r io.Reader) *bufio.Reader {
		assert.Fail(t, "don't reach here")
		return nil
	}
	t.Run("cancel before search log", func(t *testing.T) {
		fh.Close()
		os.Remove(logf)
		fh, _ = os.Create(logf)
		fh.WriteString("FATAL\nFATAL\n")

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		w, c, errLines, err := opts.searchLog(ctx, logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "should not be detected because it's not running a search")
		assert.Equal(t, int64(0), c, "should not be detected because it's not running a search")
		assert.Equal(t, "", errLines, "should not be detected because it's not running a search")
	})
	opts.testHookNewBufferedReader = nil
}

type slowReader struct {
	r io.Reader
	d time.Duration
	n int
}

func (r *slowReader) Read(p []byte) (int, error) {
	time.Sleep(r.d)
	return r.r.Read(p[:r.n])
}

func TestRunWithGlob(t *testing.T) {
	dir := t.TempDir()

	logf1 := filepath.Join(dir, "dummy1")
	fh1, _ := os.Create(logf1)
	t.Cleanup(func() {
		fh1.Close()
	})

	logf2 := filepath.Join(dir, "dummy2")
	fh2, _ := os.Create(logf2)
	t.Cleanup(func() {
		fh2.Close()
	})

	ptn := `FATAL`
	params := []string{dir, "-f", filepath.Join(dir, "dummy*"), "-p", ptn, "--check-first"}
	opts, _ := parseArgs(params)
	opts.prepare()

	t.Run("success", func(t *testing.T) {
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.OK, ckr.Status, "ckr.Status should be OK")
	})

	errorLine := "FATAL\n"
	t.Run("critical once", func(t *testing.T) {
		fh1.WriteString(errorLine)
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "ckr.Status should be CRITICAL")
	})

	t.Run("recover", func(t *testing.T) {
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.OK, ckr.Status, "ckr.Status should be OK")
	})

	t.Run("critical again", func(t *testing.T) {
		fh2.WriteString(errorLine)
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "ckr.Status should be CRITICAL")
	})
}

func TestRunWithZGlob(t *testing.T) {
	dir := t.TempDir()

	err := os.MkdirAll(filepath.Join(dir, "subdir"), 0755)
	if err != nil {
		t.Errorf("something went wrong")
	}

	logf1 := filepath.Join(dir, "dummy1")
	fh1, _ := os.Create(logf1)
	t.Cleanup(func() {
		fh1.Close()
	})

	logf2 := filepath.Join(dir, "subdir", "dummy2")
	fh2, _ := os.Create(logf2)
	t.Cleanup(func() {
		fh2.Close()
	})

	ptn := `FATAL`
	params := []string{dir, "-f", filepath.Join(dir, "**/dummy*"), "-p", ptn, "--check-first"}
	opts, _ := parseArgs(params)
	opts.prepare()

	t.Run("success", func(t *testing.T) {
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.OK, ckr.Status, "ckr.Status should be OK")
	})

	errorLine := "FATAL\n"
	t.Run("critical once", func(t *testing.T) {
		fh1.WriteString(errorLine)
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "ckr.Status should be CRITICAL")
	})

	t.Run("recover", func(t *testing.T) {
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.OK, ckr.Status, "ckr.Status should be OK")
	})

	t.Run("critical again", func(t *testing.T) {
		fh2.WriteString(errorLine)
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "ckr.Status should be CRITICAL")
	})
}

func TestRunWithMiddleOfLine(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	t.Cleanup(func() {
		fh.Close()
	})

	ptn := `FATAL`
	opts, _ := parseArgs([]string{"-s", dir, "-f", logf, "-p", ptn, "--check-first"})
	opts.prepare()

	stateFile := getStateFile(opts.StateDir, logf, opts.origArgs)

	bytes, _ := getBytesToSkip(stateFile)
	assert.Equal(t, int64(0), bytes, "should be a 0-byte indicated value")

	t.Run("middle of line", func(t *testing.T) {
		fh.WriteString("FATA")
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "it is not an error string, so it should be zero.")
		assert.Equal(t, int64(0), c, "it is not an error string, so it should be zero.")
		assert.Equal(t, "", errLines, "it is not an error string, so it should be empty.")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(0), bytes, "the pointer should be zero, because this file is not ended newline")
	})

	t.Run("fail", func(t *testing.T) {
		fh.WriteString("L\nSUCC")
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "there are one error string, so one should be detected")
		assert.Equal(t, int64(1), c, "there are one error string, so one should be detected")
		assert.Equal(t, "FATAL\n", errLines, "it should detect one line.")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len("FATAL\n")), bytes, "it should move up to the size of a single line `FATAL\n`")
	})
}

func TestRunWithNoState(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	t.Cleanup(func() {
		fh.Close()
	})

	ptn := `FATAL`
	opts, _ := parseArgs([]string{"-s", dir, "-f", logf, "-p", ptn, "--no-state"})
	opts.prepare()

	fatal := "FATAL\n"
	t.Run("two lines", func(t *testing.T) {
		fh.WriteString(fatal)
		fh.WriteString(fatal)
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(2), w, "there are two error strings, so two should be detected")
		assert.Equal(t, int64(2), c, "there are two error strings, so two should be detected")
		assert.Equal(t, strings.Repeat(fatal, 2), errLines, "it should move up to the size of a two lines `FATAL\n`")
	})

	// Make sure the entire file is loaded.
	// Because do not use the state file.
	t.Run("added one line", func(t *testing.T) {
		fh.WriteString(fatal)
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(3), w, "there are three error strings, so three should be detected")
		assert.Equal(t, int64(3), c, "there are three error strings, so three should be detected")
		assert.Equal(t, strings.Repeat(fatal, 3), errLines, "it should move up to the size of a three lines `FATAL\n`")
	})
}

func TestSearchReaderWithLevel(t *testing.T) {
	dir := t.TempDir()

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
		Pattern:         []string{`FATAL level:([0-9]+)`},
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
	warnNum, critNum, readBytes, errLines, _ := opts.searchReader(context.Background(), r)

	assert.Equal(t, int64(2), warnNum, "warnNum should be 2")
	assert.Equal(t, int64(1), critNum, "critNum should be 1")
	assert.Equal(t, "FATAL level:22\nFatal level:17\n", errLines, "invalid errLines")
	assert.Equal(t, int64(len(content)), readBytes, "readBytes should be 26")
}

func TestRunWithEncoding(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	t.Cleanup(func() {
		fh.Close()
	})

	opts, _ := parseArgs([]string{"-s", dir, "-f", logf, "-p", `エラー`, "--encoding", "euc-jp", "--check-first"})
	opts.prepare()

	t.Run("encoding", func(t *testing.T) {
		fh.Write([]byte("\xa5\xa8\xa5\xe9\xa1\xbc\n")) // エラー
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "there are one error strings, so one should be detected")
		assert.Equal(t, int64(1), c, "there are one error strings, so one should be detected")
		assert.Equal(t, "エラー\n", errLines, "given string should be detected")

		fh.Write([]byte("\xb0\xdb\xbe\xef\n")) // 異常
		w, c, errLines, err = opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(0), w, "it is not an error string, so it should be zero.")
		assert.Equal(t, int64(0), c, "it is not an error string, so it should be zero.")
		assert.Equal(t, "", errLines, "it is not an error string, so it should be empty.")

		fh.Write([]byte("\xa5\xa8\xa5\xe9\xa1\xbc\n")) // エラー
		w, c, errLines, err = opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "there are one error strings, so one should be detected")
		assert.Equal(t, int64(1), c, "there are one error strings, so one should be detected")
		assert.Equal(t, "エラー\n", errLines, "given string should be detected")
	})
}

func TestRunWithoutEncoding(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	t.Cleanup(func() {
		fh.Close()
	})

	opts, _ := parseArgs([]string{"-s", dir, "-f", logf, "-p", `エラー`, "--check-first"})
	opts.prepare()

	fatal := "\xa5\xa8\xa5\xe9\xa1\xbc\nエラー\n" // エラー
	t.Run("without encoding", func(t *testing.T) {
		fh.Write([]byte(fatal))
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "there are one error strings, so one should be detected")
		assert.Equal(t, int64(1), c, "there are one error strings, so one should be detected")
		assert.Equal(t, "エラー\n", errLines, "should be detected")
	})

	fatal = "エラー\n"
	t.Run("with encoding", func(t *testing.T) {
		fh.Write([]byte(fatal))
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "there are one error strings, so one should be detected")
		assert.Equal(t, int64(1), c, "there are one error strings, so one should be detected")
		assert.Equal(t, "エラー\n", errLines, "should be detected")
	})
}

func TestRunWithMissingOk(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")

	ptn := `FATAL`
	missing := `OK`
	params := []string{"-s", dir, "-f", logf, "-p", ptn, "--missing", missing}
	opts, _ := parseArgs(params)
	opts.prepare()

	t.Run("log file missing", func(t *testing.T) {
		ckr := run(context.Background(), params)
		assert.Equal(t, ckr.Status, checkers.OK, "ckr.Status should be OK")
		msg := fmt.Sprintf("0 warnings, 0 criticals for pattern /FATAL/.\nThe following 1 files are missing.\n%s", logf)
		assert.Equal(t, ckr.Message, msg, "no file, no error detected")
	})
}

func TestRunWithMissingWarning(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")

	ptn := `FATAL`
	missing := `WARNING`
	params := []string{"-s", dir, "-f", logf, "-p", ptn, "--missing", missing}
	opts, _ := parseArgs(params)
	opts.prepare()

	t.Run("log file missing", func(t *testing.T) {
		ckr := run(context.Background(), params)
		assert.Equal(t, ckr.Status, checkers.WARNING, "ckr.Status should be WARNING")
		msg := fmt.Sprintf("0 warnings, 0 criticals for pattern /FATAL/.\nThe following 1 files are missing.\n%s", logf)
		assert.Equal(t, ckr.Message, msg, "no file, no error detected")
	})
}

func TestRunWithMissingCritical(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")

	ptn := `FATAL`
	missing := `CRITICAL`
	params := []string{"-s", dir, "-f", logf, "-p", ptn, "--missing", missing}
	opts, _ := parseArgs(params)
	opts.prepare()

	t.Run("log file missing", func(t *testing.T) {
		ckr := run(context.Background(), params)
		assert.Equal(t, ckr.Status, checkers.CRITICAL, "ckr.Status should be CRITICAL")
		msg := fmt.Sprintf("0 warnings, 0 criticals for pattern /FATAL/.\nThe following 1 files are missing.\n%s", logf)
		assert.Equal(t, ckr.Message, msg, "no file, no error detected")
	})
}

func TestRunWithMissingUnknown(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")

	ptn := `FATAL`
	missing := `UNKNOWN`
	params := []string{"-s", dir, "-f", logf, "-p", ptn, "--missing", missing}
	opts, _ := parseArgs(params)
	opts.prepare()

	t.Run("log file missing", func(t *testing.T) {
		ckr := run(context.Background(), params)
		assert.Equal(t, ckr.Status, checkers.UNKNOWN, "ckr.Status should be UNKNOWN")
		msg := fmt.Sprintf("0 warnings, 0 criticals for pattern /FATAL/.\nThe following 1 files are missing.\n%s", logf)
		assert.Equal(t, ckr.Message, msg, "no file, no error detected")
	})
}

func TestRunWithGlobAndMissingWarning(t *testing.T) {
	dir := t.TempDir()

	logfGlob := filepath.Join(dir, "dummy*")

	ptn := `FATAL`
	missing := `WARNING`
	params := []string{"-s", dir, "-f", logfGlob, "-p", ptn, "--missing", missing}
	opts, _ := parseArgs(params)
	opts.prepare()

	t.Run("log file missing", func(t *testing.T) {
		ckr := run(context.Background(), params)
		assert.Equal(t, ckr.Status, checkers.WARNING, "ckr.Status should be WARNING")
		msg := fmt.Sprintf("0 warnings, 0 criticals for pattern /FATAL/.\nThe following 1 files are missing.\n%s", logfGlob)
		assert.Equal(t, ckr.Message, msg, "no file, no error detected")
	})
}

func TestRunMultiplePattern(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	t.Cleanup(func() {
		fh.Close()
	})

	ptn1 := `FATAL`
	ptn2 := `TESTAPPLICATION`
	params := []string{"-s", dir, "-f", logf, "-p", ptn1, "-p", ptn2}
	opts, _ := parseArgs(params)
	opts.prepare()

	stateFile := getStateFile(opts.StateDir, logf, opts.origArgs)

	bytes, _ := getBytesToSkip(stateFile)
	assert.Equal(t, int64(0), bytes, "should be a 0-byte indicated value")

	individualTextslines := "FATAL\nTESTAPPLICATION\n"
	t.Run("two lines", func(t *testing.T) {
		fh.WriteString(individualTextslines)
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.OK, ckr.Status, "ckr.Status should be OK")
		msg := "0 warnings, 0 criticals for pattern /FATAL/ and /TESTAPPLICATION/."
		assert.Equal(t, ckr.Message, msg, "it is not meet the conditions to be detected as an error string, so should be no error detected")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(individualTextslines)), bytes, "the pointer should be moved by the amount of the character string written in the file")
	})

	matchedTextsLines := "FATAL TESTAPPLICATION\nTESTAPPLICATION FATAL\n"
	t.Run("condition", func(t *testing.T) {
		fh.WriteString(matchedTextsLines)
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "ckr.Status should be CRITICAL")
		msg := "2 warnings, 2 criticals for pattern /FATAL/ and /TESTAPPLICATION/."
		assert.Equal(t, ckr.Message, msg, "it is meet the conditions to be detected as an error string, so should be error detected")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(individualTextslines)+len(matchedTextsLines)), bytes, "the pointer should be moved by the amount of the character string written in the file")
	})

	okLine := "OK\n"
	t.Run("with level", func(t *testing.T) {
		fh.WriteString(okLine)
		params := []string{"-s", dir, "-f", logf, "-p", ptn1, "-p", ptn2, "--warning-level", "12"}
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.UNKNOWN, ckr.Status, "ckr.Status should be UNKNOWN")
		msg := "When multiple patterns specified, --warning-level --critical-level can not be used"
		assert.Equal(t, ckr.Message, msg, "it should be detected that the option cannot be specified")
	})

	t.Run("invalid pattern", func(t *testing.T) {
		fh.WriteString(okLine)
		ptn3 := "+"
		params := []string{"-s", dir, "-f", logf, "-p", ptn1, "-p", ptn3}
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.UNKNOWN, ckr.Status, "ckr.Status should be UNKNOWN")
		msg := "pattern is invalid"
		assert.Equal(t, ckr.Message, msg, "it should be detected that the pattern cannot be specified")
	})
}

func TestRunWithSuppressOption(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	t.Cleanup(func() {
		fh.Close()
	})

	ptn1 := `FATAL`
	ptn2 := `TESTAPPLICATION`
	params := []string{"-s", dir, "-f", logf, "-p", ptn1, "-p", ptn2, "--suppress-pattern"}
	opts, _ := parseArgs(params)
	opts.prepare()

	stateFile := getStateFile(opts.StateDir, logf, opts.origArgs)

	bytes, _ := getBytesToSkip(stateFile)
	assert.Equal(t, int64(0), bytes, "should be a 0-byte indicated value")

	individualTextslines := "FATAL\nTESTAPPLICATION\n"
	t.Run("two lines", func(t *testing.T) {
		fh.WriteString(individualTextslines)
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.OK, ckr.Status, "ckr.Status should be OK")
		msg := "0 warnings, 0 criticals."
		assert.Equal(t, ckr.Message, msg, "it is not meet the conditions to be detected as an error string, so should be no error detected")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(individualTextslines)), bytes, "the pointer should be moved by the amount of the character string written in the file")
	})

	matchedTextsLines := "FATAL TESTAPPLICATION\nTESTAPPLICATION FATAL\n"
	t.Run("condition", func(t *testing.T) {
		fh.WriteString(matchedTextsLines)
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.CRITICAL, ckr.Status, "ckr.Status should be CRITICAL")
		msg := "2 warnings, 2 criticals."
		assert.Equal(t, ckr.Message, msg, "it is meet the conditions to be detected as an error string, so should be error detected")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(individualTextslines)+len(matchedTextsLines)), bytes, "the pointer should be moved by the amount of the character string written in the file")
	})

	okLine := "OK\n"
	t.Run("with level", func(t *testing.T) {
		fh.WriteString(okLine)
		params := []string{"-s", dir, "-f", logf, "-p", ptn1, "-p", ptn2, "--warning-level", "12", "--suppress-pattern"}
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.UNKNOWN, ckr.Status, "ckr.Status should be UNKNOWN")
		msg := "When multiple patterns specified, --warning-level --critical-level can not be used"
		assert.Equal(t, ckr.Message, msg, "it should be detected that the option cannot be specified")
	})

	t.Run("invalid pattern", func(t *testing.T) {
		fh.WriteString(okLine)
		ptn3 := "+"
		params := []string{"-s", dir, "-f", logf, "-p", ptn1, "-p", ptn3, "--suppress-pattern"}
		ckr := run(context.Background(), params)
		assert.Equal(t, checkers.UNKNOWN, ckr.Status, "ckr.Status should be UNKNOWN")
		msg := "pattern is invalid"
		assert.Equal(t, ckr.Message, msg, "it should be detected that the pattern cannot be specified")
	})
}

func TestRunMultipleExcludePattern(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	t.Cleanup(func() {
		fh.Close()
	})

	params := []string{"-s", dir, "-f", logf, "-p", "ERROR", "-p", "TESTAPP", "-E", "FOO", "-E", "BAR"}
	opts, _ := parseArgs(params)
	opts.prepare()

	stateFile := getStateFile(opts.StateDir, logf, opts.origArgs)

	bytes, _ := getBytesToSkip(stateFile)
	assert.Equal(t, int64(0), bytes, "stateFile size should be 0 (actual: %d)", bytes)

	tests := []struct {
		logmsg       string
		want         checkers.Status
		matchedCount int64
	}{
		{
			logmsg:       "[TESTAPP] DEBUG: THIS LINE UNMATCHED\n",
			want:         checkers.OK,
			matchedCount: 0,
		},
		{
			logmsg:       "[TESTAPP] ERROR: THIS LINE MATCHED\n",
			want:         checkers.CRITICAL,
			matchedCount: 1,
		},
		{
			logmsg:       "[TESTAPP] ERROR: THIS LINE EXCLUDED(FOO BAR)\n[TESTAPP] ERROR: THIS LINE MATCHED(FOO)\n",
			want:         checkers.CRITICAL,
			matchedCount: 1,
		},
		{
			logmsg:       "[TESTAPP] ERROR: THIS LINE EXCLUDED(FOO BAR)\nERROR: THIS LINE UNMATCHED\n",
			want:         checkers.OK,
			matchedCount: 0,
		},
	}

	for _, test := range tests {
		fh.WriteString(test.logmsg)
		ckr := run(context.Background(), params)

		assert.Equal(t, test.want, ckr.Status, fmt.Sprintf("ckr.Status should be %s", test.want.String()))

		msg := fmt.Sprintf("%d warnings, %d criticals for pattern /ERROR/ and /TESTAPP/.", test.matchedCount, test.matchedCount)
		assert.Equal(t, msg, ckr.Message, fmt.Sprintf("chk.Message should be '%s' (actual: '%s')", msg, ckr.Message))
	}
}

func TestParseFilePattern(t *testing.T) {
	dir := t.TempDir()

	logf1 := filepath.Join(dir, "dummy1.txt")
	fh1, _ := os.Create(logf1)
	fh1.Close()

	logf2 := filepath.Join(dir, "dummy2.txt")
	fh2, _ := os.Create(logf2)
	fh2.Close()

	logf3 := filepath.Join(dir, "DUMMY3.txt")
	fh3, _ := os.Create(logf3)
	fh3.Close()

	tests := []struct {
		desc        string
		directory   string
		filePattern string
		insensitive bool
		actual      []string
		skipWindows bool
	}{
		{
			desc:        "filePattern only",
			directory:   "",
			filePattern: dir + string(filepath.Separator) + `dummy\d.txt`,
			insensitive: true,
			actual:      []string{logf1, logf2, logf3},
			// If in Windows, '\\' is used as both path separator and special meaning of regular expression.
			// Thus the result of filepath.Dir do not become a expected directory name.
			skipWindows: true,
		},
		{
			desc:        "separate directory, file pattern",
			directory:   dir,
			filePattern: `dummy\d.txt`,
			insensitive: true,
			actual:      []string{logf1, logf2, logf3},
			skipWindows: false,
		},
		{
			desc:        "case sensitive",
			directory:   dir,
			filePattern: `DUMMY\d.txt`,
			insensitive: false,
			actual:      []string{logf3},
			skipWindows: false,
		},
	}

	for _, tt := range tests {
		if tt.skipWindows && runtime.GOOS == "windows" {
			continue
		}
		list, err := parseFilePattern(tt.directory, tt.filePattern, tt.insensitive)
		if err != nil {
			t.Fatalf("failed: %s", err)
		}
		sort.Strings(list)
		sort.Strings(tt.actual)
		assert.Equal(t, list, tt.actual, fmt.Sprintf("%s - list should be '%s' (actual: '%s')", tt.desc, list, tt.actual))
	}
}
