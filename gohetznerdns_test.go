package gohetznerdns

import (
	"testing"

	"gotest.tools/assert"
)

func TestNewClient(t *testing.T) {
	dns, _ := NewClient("asdadasdasd")
	assert.Assert(t, dns != nil)
	assert.Assert(t, dns.client != nil)
	assert.Assert(t, dns.ZoneService != nil)
	assert.Assert(t, dns.RecordService != nil)
}

func TestNewClientTokenError(t *testing.T) {
	_, err := NewClient("   ")
	assert.Error(t, err, "token is invalid because cannot be empty")
}

func TestSetBaseUrl(t *testing.T) {
	dns, _ := NewClient("asdaasdada")
	err := dns.SetBaseURL("https://test.com")
	assert.NilError(t, err)
}
func TestSetBaseUrlInvalid(t *testing.T) {
	dns, _ := NewClient("asdasdsadas")
	err := dns.SetBaseURL("https://te|.^com")
	assert.Error(t, err, "parse \"https://te|.^com\": invalid character \"|\" in host name")
}
