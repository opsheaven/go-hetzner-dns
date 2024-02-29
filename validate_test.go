package gohetznerdns

import (
	"testing"

	"gotest.tools/assert"
)

func TestArgError(t *testing.T) {
	expected := "foo is invalid because bar"
	err := newArgError("foo", "bar")
	assert.Equal(t, err.Error(), expected)
}

func TestValidateNotEmpty(t *testing.T) {
	err := validateNotEmpty("test", "test")
	assert.NilError(t, err)
}

func TestValidateEmpty(t *testing.T) {
	err := validateNotEmpty("test", "    ")
	assert.Error(t, err, "test is invalid because cannot be empty")
}
