package main

import (
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestEscapedString(t *testing.T) {
	assert.Equal(t, "\n", escapedString(`\n`), "something went wrong")
	assert.Equal(t, "hoge\\", escapedString(`hoge\`), "something went wrong")
	assert.Equal(t, "ho\rge", escapedString(`ho\rge`), "something went wrong")
	assert.Equal(t, "ho\\oge", escapedString(`ho\oge`), "something went wrong")
	assert.Equal(t, "", escapedString(``), "something went wrong")
}

func TestTLS(t *testing.T) {
	opts, err := parseArgs([]string{"-S", "-H", "www.verisign.com", "-p", "443"})
	assert.Equal(t, nil, err, "no errors")
	ckr := opts.run()
	assert.Equal(t, checkers.OK, ckr.Status, "should be OK")
}

func TestFTP(t *testing.T) {
	opts, err := parseArgs([]string{"--service=ftp", "-H", "ftp.iij.ad.jp"})
	assert.Equal(t, nil, err, "no errors")
	ckr := opts.run()
	assert.Equal(t, checkers.OK, ckr.Status, "should be OK")
}
