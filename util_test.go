package gohetznerdns

import (
	"testing"

	"gotest.tools/assert"
)

func TestIIFValid(t *testing.T) {
	assert.Equal(t, iif(1 == 1, 1, 2), 1)
}
func TestIIFInValid(t *testing.T) {
	assert.Equal(t, iif(1 == 2, 1, 2), 2)
}

func TestValidateNotNil(t *testing.T) {
	d := "data"
	assert.NilError(t, validateNotNil("test", &d))
}
func TestValidateNotEmpty(t *testing.T) {
	d := "data"
	assert.NilError(t, validateNotEmpty("test", &d))
}

func TestValidateNotEmptyWhenNil(t *testing.T) {
	var d *string
	assert.Error(t, validateNotEmpty("test", d), "900 : test is nil")
}
