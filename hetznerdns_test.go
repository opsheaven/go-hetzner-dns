package hetznerdns

import (
	"testing"

	"gotest.tools/assert"
)

func TestNewClient(t *testing.T) {
	dns := NewClient()
	assert.Assert(t, dns != nil)
	assert.Assert(t, dns.client != nil)
	assert.Assert(t, dns.ZoneService != nil)
	assert.Assert(t, dns.RecordService != nil)
}

func TestSetBaseUrl(t *testing.T) {
	err := NewClient().SetBaseURL("https://test.com")
	assert.NilError(t, err)
}
func TestSetBaseUrlInvalid(t *testing.T) {
	err := NewClient().SetBaseURL("https://te|.^com")
	assert.Error(t, err, "parse \"https://te|.^com\": invalid character \"|\" in host name")
}
func TestSetToken(t *testing.T) {
	err := NewClient().SetToken("asdasdsadsadda")
	assert.NilError(t, err)
}

func TestSetTokenEmpty(t *testing.T) {
	err := NewClient().SetToken("    ")
	assert.Error(t, err, "token is invalid because cannot be empty")
}
