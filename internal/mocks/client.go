package mocks

import (
	"io"
	"net/http"
)

// MockClient is the mock client
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
	req    *http.Request
}

var (
	// GetDoFunc fetches the mock client's `Do` func
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

// Do is the mock client's `Do` func
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	m.req = req
	return GetDoFunc(req)
}

func (m *MockClient) GetRequestBody() io.ReadCloser {
	return m.req.Body
}
