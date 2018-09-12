// +build !windows

package checklog

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFileByInode(t *testing.T) {
	dir, err := ioutil.TempDir("", "check-log-test")
	if err != nil {
		t.Errorf("something went wrong")
	}
	defer os.RemoveAll(dir)

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	defer fh.Close()

	testFileExist := func() {
		logfi, err := os.Stat(logf)
		inode := detectInode(logfi)
		f, err := findFileByInode(inode, dir)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, logf, f)
	}
	testFileExist()

	testFileNotExist := func() {
		f, err := findFileByInode(100, dir)
		assert.Equal(t, err, errFileNotFoundByInode, "err should be errFileNotFoundByInode")
		assert.Equal(t, "", f)
	}
	testFileNotExist()
}

func TestRunTraceInode(t *testing.T) {
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

	l1 := "SUCCESS\n"
	l2 := "FATAL\n"
	l3 := "SUCCESS\n"
	testRotate := func() {
		// first check
		fh.WriteString(l1)
		opts.searchLog(logf)

		// write FATAL
		fh.WriteString(l2)
		fh.Close()

		// logrotate
		rotatedLogf := filepath.Join(dir, "dummy.1")
		os.Rename(logf, rotatedLogf)
		fh, _ = os.Create(logf)

		fh.WriteString(l3)
		// second check
		w, c, errLines, err := opts.searchLog(logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "something went wrong")
		assert.Equal(t, int64(1), c, "something went wrong")
		assert.Equal(t, "FATAL\n", errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(l3)), bytes, "should not include oldfile skip bytes")

		os.Remove(rotatedLogf)
	}
	testRotate()
}
