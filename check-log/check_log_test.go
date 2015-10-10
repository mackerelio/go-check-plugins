package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStateFile(t *testing.T) {
	sPath := getStateFile("/var/lib", "C:/Windows/hoge")
	assert.Equal(t, sPath, "/var/lib/C/Windows/hoge", "drive letter should be cared")

	sPath = getStateFile("/var/lib", "/linux/hoge")
	assert.Equal(t, sPath, "/var/lib/linux/hoge", "drive letter should be cared")
}
