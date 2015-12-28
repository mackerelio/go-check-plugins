package main

import (
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	ckr := run([]string{"-u", "hoge"})
	assert.Equal(t, ckr.Status, checkers.CRITICAL, "chr.Status should be CRITICAL")
	assert.Equal(t, ckr.Message, `Get hoge: unsupported protocol scheme ""`, "something went wrong")
}

func TestNoCheckCertificate(t *testing.T) {
	ckr := run([]string{"--no-check-certificate", "-u", "hoge"})
	assert.Equal(t, ckr.Status, checkers.CRITICAL, "chr.Status should be CRITICAL")
	assert.Equal(t, ckr.Message, `Get hoge: unsupported protocol scheme ""`, "something went wrong")
}
