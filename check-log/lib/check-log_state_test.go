package checklog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Format implements fmt.Formatter.
func (s *state) Format(f fmt.State, c rune) {
	if s == nil {
		fmt.Fprintf(f, "<nil>")
		return
	}
	fmt.Printf("%"+string(c), *s)
}

func TestLoadStateIfFileNotExist(t *testing.T) {
	file := "testdata/file_not_found"
	s, err := loadState(file)
	if err != nil {
		t.Errorf("loadState(%q) = %v; want nil", file, err)
	}
	if s != nil {
		t.Errorf("loadState(%q) = %v; want nil", file, *s)
	}
}

func TestLoadStateIfAccessDenied(t *testing.T) {
	var file string
	switch runtime.GOOS {
	case "windows":
		file = `C:\pagefile.sys`
	default:
		file = "testdata/plain.txt/file"
	}
	s, err := loadState(file)
	if err == nil {
		t.Errorf("loadState(%q) = %v; want an error", file, s)
	}
}

// TODO(lufia): We might be better to test a condition too that file is exist but loadState can't read it.

func TestSaveStateIfFileNotExist(t *testing.T) {
	file := "testdata/file_will_create"
	defer func() {
		err := os.Remove(file)
		if err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}()

	s := &state{
		SkipBytes: 15,
		Inode:     150,
	}
	testSaveLoadState(t, file, s)
}

func TestSaveStateOverwrittenIfFileExist(t *testing.T) {
	file := "testdata/state_overwritten"
	defer func() {
		err := os.Remove(file)
		if err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}()

	err := ioutil.WriteFile(file, []byte(`{"skip_bytes": 10, "inode": 100}`), 0644)
	if err != nil {
		t.Errorf("WriteFile: %v", err)
		return
	}
	s := &state{
		SkipBytes: 15,
		Inode:     150,
	}
	testSaveLoadState(t, file, s)
}

func testSaveLoadState(t *testing.T, file string, s *state) {
	t.Helper()

	if err := saveState(file, s); err != nil {
		t.Errorf("saveState(%v) = %v; want nil", file, *s)
		return
	}
	s1, err := loadState(file)
	if err != nil {
		t.Errorf("loadState: %v", err)
		return
	}
	if !reflect.DeepEqual(s, s1) {
		t.Errorf("saveState(%v) -> loadState() = %v", s, s1)
	}
}

func TestSaveStateIfAccessDenied(t *testing.T) {
	switch runtime.GOOS {
	case "windows": // Almost os.Chmod operations will ignore on Windows.
		t.Skip()
		return
	}
	file := "testdata/readonly/state"
	dir := filepath.Dir(file)
	defer func() {
		if err := os.Chmod(dir, 0755); err != nil {
			t.Fatalf("Chmod: %v", err)
		}
		err := os.RemoveAll(dir)
		if err != nil && !os.IsNotExist(err) {
			t.Fatalf("RemoveAll: %v", err)
		}
	}()

	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Errorf("MkdirAll: %v", err)
		return
	}
	data := []byte(`{"skip_bytes": 10, "inode": 100}`)
	if err := ioutil.WriteFile(file, data, 0644); err != nil {
		t.Errorf("WriteFile: %v", err)
		return
	}
	if err := os.Chmod(dir, 0500); err != nil {
		t.Errorf("Chmod: %v", err)
		return
	}
	s := &state{
		SkipBytes: 15,
		Inode:     150,
	}
	saveState(file, s) // an error can be ignored in this case.
	data1, err := ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("ReadFile: %v", err)
		return
	}
	if !bytes.Equal(data1, data) {
		t.Errorf("saveState into readonly directory should keep original contents: result = %s", data1)
	}
}

func TestGetBytesToSkipOld(t *testing.T) {
	file := "testdata/old-state.txt"
	n, err := getBytesToSkipOld(file)
	if err != nil {
		t.Fatal(err)
	}
	if want := int64(20); n != want {
		t.Errorf("getBytesToSkip(%s) = %d; want %d", file, n, want)
	}
}

func TestGetBytesToSkipOldIfFileNotExist(t *testing.T) {
	file := "testdata/file_not_found"
	n, err := getBytesToSkipOld(file)
	if err != nil {
		t.Fatal(err)
	}
	if want := int64(0); n != want {
		t.Errorf("getBytesToSkip(%s) = %d; want %d", file, n, want)
	}
}

func TestGetBytesToSkipOldErr(t *testing.T) {
	var file string
	switch runtime.GOOS {
	case "windows":
		file = `C:\pagefile.sys`
	default:
		file = "testdata/plain.txt/file"
	}
	n, err := getBytesToSkipOld(file)
	if err == nil {
		t.Errorf("getBytesToSkip(%s) = %d; but want an error", file, n)
	}
}
