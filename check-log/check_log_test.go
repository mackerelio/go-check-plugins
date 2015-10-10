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

func TestWriteBytesToSkip(t *testing.T) {
	f := ".tmp/fuga/piyo"
	err := writeBytesToSkip(f, 15)
	assert.Equal(t, err, nil, "err should be nil")

	skipBytes, err := getBytesToSkip(f)
	assert.Equal(t, err, nil, "err should be nil")
	assert.Equal(t, skipBytes, int64(15))
}
