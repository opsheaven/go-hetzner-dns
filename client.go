package gohetznerdns

import (
	"encoding/json"
	"fmt"
	"net/url"
	"slices"

	"github.com/go-resty/resty/v2"
)

const (
	apiVersion            = "1.1.1"
	defaultBaseURL        = "https://dns.hetzner.com"
	basePath              = "/api/v1"
	contentTypeJson       = "application/json; charset=utf-8"
	contentTypeText       = "text/plain"
	contentTypeURLEncoded = "application/x-www-form-urlencoded; charset=utf-8"
)

type client struct {
	client  *resty.Client
	baseURL *url.URL
	token   string
}

type request struct {
	request             *resty.Request
	baseURL             *url.URL
	expectedStatusCodes []int
	result              interface{}
}

func newClient() *client {
	client := &client{client: resty.New()}
	client.setBaseURL(defaultBaseURL)
	return client
}

func (client *client) setBaseURL(baseUrl string) error {
	baseURL, err := url.Parse(baseUrl)
	if err == nil {
		client.baseURL = baseURL
	}
	return err
}

func (client *client) setToken(token string) {
	client.token = token
}

func (c *client) createRequest(contentType string, expectedStatusCodes ...int) *request {
	request := &request{
		request:             c.client.R(),
		baseURL:             c.baseURL,
		expectedStatusCodes: expectedStatusCodes,
	}
	request.request.SetHeader("Content-Type", contentType).
		SetHeader("Auth-API-Token", c.token)
	return request
}

func (c *client) createJsonRequest(expectedStatusCodes ...int) *request {
	return c.createRequest(contentTypeJson, expectedStatusCodes...)
}

func (c *client) createTextRequest(expectedStatusCodes ...int) *request {
	return c.createRequest(contentTypeText, expectedStatusCodes...)
}

func (r *request) setQueryParams(params map[string]string) *request {
	r.request.SetQueryParams(params)
	return r
}

func (r *request) setBody(body interface{}) *request {
	r.request.SetBody(body)
	return r
}

func (r *request) setResult(result interface{}) *request {
	r.result = result
	return r
}

func (r *request) execute(method, path string) ([]byte, error) {
	var u *url.URL
	var err error
	if u, err = r.baseURL.Parse(basePath + path); err != nil {
		return nil, err
	}
	response, err := r.request.Execute(method, u.String())
	if err != nil {
		return nil, err
	}
	body := response.Body()
	if r.result != nil && body != nil {
		if err := json.Unmarshal(response.Body(), r.result); err != nil {
			return body, err
		}
	}
	if slices.Contains(r.expectedStatusCodes, response.StatusCode()) {
		return body, nil
	}

	return body, fmt.Errorf(response.Status())
}
