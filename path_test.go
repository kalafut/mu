package jpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitExt(t *testing.T) {
	assert := assert.New(t)

	base, ext := SplitExt("foo.bar")
	assert.Equal(base, "foo")
	assert.Equal(ext, ".bar")

	base, ext = SplitExt("a/b/c/foo.bar")
	assert.Equal(base, "a/b/c/foo")
	assert.Equal(ext, ".bar")

	base, ext = SplitExt("foo")
	assert.Equal(base, "foo")
	assert.Equal(ext, "")

	base, ext = SplitExt(".bar")
	assert.Equal(base, "")
	assert.Equal(ext, ".bar")
}
