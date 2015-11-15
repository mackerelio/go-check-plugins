package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapedString(t *testing.T) {
	assert.Equal(t, "\n", escapedString(`\n`), "something went wrong")
	assert.Equal(t, "hoge\\", escapedString(`hoge\`), "something went wrong")
	assert.Equal(t, "ho\rge", escapedString(`ho\rge`), "something went wrong")
	assert.Equal(t, "ho\\oge", escapedString(`ho\oge`), "something went wrong")
	assert.Equal(t, "", escapedString(``), "something went wrong")
}
