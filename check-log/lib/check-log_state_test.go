package checklog

import (
	"fmt"
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

func TestLoadStateIfFileNotFound(t *testing.T) {
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
	file := "testdata/file.txt/any"
	s, err := loadState(file)
	if err == nil {
		t.Errorf("loadState(%q) = %v; want an error", file, s)
	}
}
