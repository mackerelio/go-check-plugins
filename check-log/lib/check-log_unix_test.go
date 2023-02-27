//go:build !windows

package checklog

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFileByInode(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	defer fh.Close()

	testFileExist := func() {
		logfi, err := os.Stat(logf)
		assert.Equal(t, err, nil, "err should be nil")
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

func TestOpenOldFile(t *testing.T) {
	dir := t.TempDir()

	logf := filepath.Join(dir, "dummy")
	fh, _ := os.Create(logf)
	defer fh.Close()

	testFoundOldFile := func() {
		ologf := filepath.Join(dir, "dummy.1")
		ofh, _ := os.Create(ologf)
		ofi, _ := ofh.Stat()
		defer ofh.Close()

		l1 := "FATAL\n"
		ofh.WriteString(l1 + l1)

		state := &state{SkipBytes: int64(len(l1)), Inode: detectInode(ofi)}
		f, err := openOldFile(logf, state)
		defer f.Close()
		pos, _ := f.Seek(-1, io.SeekCurrent) // get current offset

		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, ofh.Name(), f.Name(), "openOldFile should be return old file")
		assert.Equal(t, len(l1)-1, int(pos))
	}
	testFoundOldFile()

	testNotFoundOldFile := func() {
		state := &state{SkipBytes: 0, Inode: 0}
		f, err := openOldFile(logf, state)

		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, f, (*os.File)(nil), "file should be nil")
	}
	testNotFoundOldFile()

	testFoundOldFileByInode := func() {
		// dir different from logf
		dir := t.TempDir()
		ologf := filepath.Join(dir, "dummy.1")

		ofh, _ := os.Create(ologf)
		ofi, _ := ofh.Stat()
		defer ofh.Close()

		state := &state{SkipBytes: 0, Inode: detectInode(ofi)}
		f, err := openOldFile(logf, state)
		defer f.Close()

		assert.Equal(t, err, nil, "err should be nil")
	}
	testFoundOldFileByInode()
}

func TestRunTraceInode(t *testing.T) {
	dir := t.TempDir()

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
		opts.searchLog(context.Background(), logf)

		// write FATAL
		fh.WriteString(l2)
		fh.Close()

		// logrotate
		rotatedLogf := filepath.Join(dir, "dummy.1")
		os.Rename(logf, rotatedLogf)
		fh, _ = os.Create(logf)

		fh.WriteString(l3)
		// second check
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil")
		assert.Equal(t, int64(1), w, "something went wrong")
		assert.Equal(t, int64(1), c, "something went wrong")
		assert.Equal(t, "FATAL\n", errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(l3)), bytes, "should not include oldfile skip bytes")

		os.Remove(rotatedLogf)
	}
	testRotate()

	// case of renaming to different directory
	// from the directory of the target logfile
	testRotateDifferentDir := func() {
		// first check
		fh.WriteString(l1)
		opts.searchLog(context.Background(), logf)

		// write FATAL
		fh.WriteString(l2)
		fh.Close()

		// logrotate to <dir>/dummy/dummy.1
		newDir := filepath.Join(dir, "dummy")
		os.MkdirAll(newDir, 0755)
		rotatedLogf := filepath.Join(newDir, "dummy.1")
		os.Rename(logf, rotatedLogf)
		fh, _ = os.Create(logf)

		fh.WriteString(l3)
		// second check
		w, c, errLines, err := opts.searchLog(context.Background(), logf)
		assert.Equal(t, err, nil, "err should be nil when the old file is not found")
		assert.Equal(t, int64(0), w, "something went wrong")
		assert.Equal(t, int64(0), c, "something went wrong")
		assert.Equal(t, "", errLines, "something went wrong")

		bytes, _ = getBytesToSkip(stateFile)
		assert.Equal(t, int64(len(l3)), bytes, "should not include oldfile skip bytes")

		os.Remove(rotatedLogf)
	}
	testRotateDifferentDir()
}
