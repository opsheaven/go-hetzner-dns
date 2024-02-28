package hetznerdns

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
)

func TestNewDNSClient(t *testing.T) {
	client := newClient()
	assert.Equal(t, client.baseURL.String(), defaultBaseURL)
}

func TestNewBaseURL(t *testing.T) {
	url := "https://test.local.com"
	client := newClient()
	client.setBaseURL(url)
	assert.Equal(t, client.baseURL.String(), url)
}

func TestNewBaseURLInvalid(t *testing.T) {
	url := "https://|^%test.local.com"
	client := newClient()
	err := client.setBaseURL(url)
	assert.Equal(t, client.baseURL.String(), defaultBaseURL)
	assert.Error(t, err, "parse \"https://|^%test.local.com\": invalid character \"|\" in host name")
}

func TestToken(t *testing.T) {
	token := "1231311"
	client := newClient()
	client.setToken(token)
	if client.token != token {
		t.Error("Client should use provided token")
	}
}

func TestCreateRequest(t *testing.T) {
	token := "1231311"
	url := "https://test.local.com"
	client := newClient()
	client.setBaseURL(url)
	client.setToken(token)
	r := client.createRequest(contentTypeText, 200)

	if r == nil {
		t.Error("Client should create a request")
	}
	if r.baseURL.String() != url {
		t.Error("Client should use configured domain")
	}
	if len(r.expectedStatusCodes) != 1 && r.expectedStatusCodes[0] != 200 {
		t.Error("Client status codes are broken")
	}
}

func TestCreateJsonRequest(t *testing.T) {
	token := "1231311"
	url := "https://test.local.com"
	client := newClient()
	client.setBaseURL(url)
	client.setToken(token)
	r := client.createJsonRequest(200)

	if r == nil {
		t.Error("Client should create a request")
	}
	if r.baseURL.String() != url {
		t.Error("Client should use configured domain")
	}
	if len(r.expectedStatusCodes) != 1 && r.expectedStatusCodes[0] != 200 {
		t.Error("Client status codes are broken")
	}
}

func TestCreateTextRequest(t *testing.T) {
	token := "1231311"
	url := "https://test.local.com"
	client := newClient()
	client.setBaseURL(url)
	client.setToken(token)
	r := client.createTextRequest(200)

	if r == nil {
		t.Error("Client should create a request")
	}
	if r.baseURL.String() != url {
		t.Error("Client should use configured domain")
	}
	if len(r.expectedStatusCodes) != 1 && r.expectedStatusCodes[0] != 200 {
		t.Error("Client status codes are broken")
	}
}

type testData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func TestClient(t *testing.T) {
	response := "{\"code\":122,\"message\":\"Message\"}"
	_body := "body"
	token := "abcabcabvc"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/test", func(w http.ResponseWriter, r *http.Request) {
		response := response
		assert.Equal(t, r.Header.Get("Auth-API-Token"), token)
		fmt.Fprint(w, response)
	})

	test := new(testData)
	client := newClient()
	client.setBaseURL(server.URL)
	client.setToken(token)
	request := client.createJsonRequest(200)

	body, err := request.
		setQueryParams(map[string]string{"a": "b"}).
		setResult(test).
		setBody(_body).
		execute("GET", "/test")

	assert.NilError(t, err)
	assert.Equal(t, string(body), response)
	assert.Equal(t, test.Code, 122)
	assert.Equal(t, test.Message, "Message")
}

func TestClientWithUnexpectedStatusses(t *testing.T) {
	response := "{\"code\":122,\"message\":\"Message\"}"
	_body := "body"
	token := "abcabcabvc"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/test", func(w http.ResponseWriter, r *http.Request) {
		response := response
		assert.Equal(t, r.Header.Get("Auth-API-Token"), token)
		fmt.Fprint(w, response)
	})

	test := new(testData)
	client := newClient()
	client.setBaseURL(server.URL)
	client.setToken(token)
	request := client.createJsonRequest(201)

	_, err := request.
		setQueryParams(map[string]string{"a": "b"}).
		setResult(test).
		setBody(_body).
		execute("GET", "/test")

	assert.Error(t, err, "200 OK")

}
