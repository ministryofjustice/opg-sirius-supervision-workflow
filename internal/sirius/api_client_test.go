package sirius

import (
	"context"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"log/slog"
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

func SetUpTest() (*slog.Logger, *mocks.MockClient) {
	logger := slog.New(slog.
		NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == "level" {
					return slog.Attr{}
				}

				if a.Key == "time" {
					a.Key = "timestamp"
				}

				if a.Key == "msg" {
					a.Key = "message"
				}

				return a
			},
		}).
		WithAttrs([]slog.Attr{slog.String("service_name", "opg-sirius-workflow")}))
	mockClient := &mocks.MockClient{}
	return logger, mockClient
}
