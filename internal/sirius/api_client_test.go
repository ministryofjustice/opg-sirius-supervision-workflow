package sirius

import (
	"context"
	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getContext(cookies []*http.Cookie) Context {
	return Context{
		Context:   context.Background(),
		Cookies:   cookies,
		XSRFToken: "abcde",
	}
}

func TestClientError(t *testing.T) {
	assert.Equal(t, "message", ClientError("message").Error())
}

func TestValidationError(t *testing.T) {
	assert.Equal(t, "message", ValidationError{Message: "message"}.Error())
}

func TestStatusError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/some/url", nil)

	resp := &http.Response{
		StatusCode: http.StatusTeapot,
		Request:    req,
	}

	err := newStatusError(resp)

	assert.Equal(t, "POST /some/url returned 418", err.Error())
	assert.Equal(t, "unexpected response from Sirius", err.Title())
	assert.Equal(t, err, err.Data())
}

func SetUpTest() (*logging.Logger, *mocks.MockClient) {
	logger := logging.New(os.Stdout, "opg-sirius-workflow ")
	mockClient := &mocks.MockClient{}
	return logger, mockClient
}
