package checkfilesize

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestSizeValue(t *testing.T) {
	// map[input] = expected float value
	var testData = map[string]float64{
		"100":  100,
		"1.2":  1.2,
		"1.2k": 1.2 * 1024,
		"1.2K": 1.2 * 1024,
		"1.2m": 1.2 * 1024 * 1024,
		"1.2M": 1.2 * 1024 * 1024,
		"1.2g": 1.2 * 1024 * 1024 * 1024,
		"1.2G": 1.2 * 1024 * 1024 * 1024,
		"1.2t": 1.2 * 1024 * 1024 * 1024 * 1024,
		"1.2T": 1.2 * 1024 * 1024 * 1024 * 1024,
		"1T":   1.0 * 1024 * 1024 * 1024 * 1024,
	}
	for input, expect := range testData {
		size, err := sizeValue(input)
		if err != nil {
			t.Error(err)
		}
		if size != expect {
			msg := fmt.Errorf("size doesn't match: input = %s, expect = %f, actual = %f", input, expect, size)
			t.Error(msg)
		}
	}
}

func TestSizeValueWithInvalidInput(t *testing.T) {
	var testData = []string{
		"1..2",
		".2k",
		"1.2aaaaa",
		"aaaaa",
	}
	for _, input := range testData {
		size, err := sizeValue(input)
		if err == nil {
			t.Error("Error should occur")
		}
		if size != -1.0 {
			t.Error("size should be -1")
		}
	}
}

func TestListFiles(t *testing.T) {
	// map[depth] = expected files
	var testData = map[int][]string{
		1: []string{"test_dir/file1"},
		2: []string{"test_dir/file1", "test_dir/depth1/file2"},
		3: []string{"test_dir/file1", "test_dir/depth1/file2", "test_dir/depth1/depth2/file3"},
		4: []string{"test_dir/file1", "test_dir/depth1/file2", "test_dir/depth1/depth2/file3"},
	}
	for depth, expects := range testData {
		files, err := listFiles("test_dir", depth)
		if err != nil {
			t.Error(err)
		}
		sort.Strings(files)
		sort.Strings(expects)
		if reflect.DeepEqual(files, expects) != true {
			msg := fmt.Errorf("file doesn't match: depth = %d, expects = %v, actual = %v", depth, expects, files)
			t.Error(msg)
		}
	}

}

func TestListFilesWithInvalidDir(t *testing.T) {
	files, err := listFiles("dir_does_not_exist", 1)
	if err == nil {
		t.Error("Error should occur")
	}
	if len(files) != 0 {
		t.Error("files should be empty")
	}
}
