package mocks

import (
	"github.com/ministryofjustice/opg-go-common/logging"
	"net/http"
	"os"
)

// MockClient is the mock client
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	// GetDoFunc fetches the mock client's `Do` func
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

// Do is the mock client's `Do` func
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func (m *MockClient) logRequest(r *http.Request, err error) {
	logger := logging.New(os.Stdout, "opg-sirius-workflow ")
	logger.Print(r.Method)
	logger.Print(r.URL.Path)
	if err != nil {
		logger.Print(err)
	}
}
