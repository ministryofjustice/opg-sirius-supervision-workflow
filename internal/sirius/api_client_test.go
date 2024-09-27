package sirius

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
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
	logger := telemetry.NewLogger("opg-sirius-workflow")
	mockClient := &mocks.MockClient{}
	return logger, mockClient
}

type LogEntry struct {
	Level   string `json:"level"`
	Message string `json:"msg"`
}

func TestLogResponseLogsErrors(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	client := ApiClient{
		logger: logger,
	}

	client.logResponse(
		&http.Request{Method: "POST", URL: &url.URL{Path: "/my-page"}},
		&http.Response{StatusCode: 200},
		errors.New("some error"),
	)

	logLines := strings.Split(strings.Trim(buf.String(), "\n"), "\n")
	assert.Equal(t, 2, len(logLines))

	logEntries := make([]LogEntry, len(logLines))
	for i, line := range logLines {
		err := json.Unmarshal([]byte(line), &logEntries[i])
		assert.Nil(t, err)
	}

	assert.Equal(t, "INFO", logEntries[0].Level)
	assert.Equal(t, "method: POST, url: /my-page, response: 200", logEntries[0].Message)

	assert.Equal(t, "ERROR", logEntries[1].Level)
	assert.Equal(t, "some error", logEntries[1].Message)
}

func TestLogResponseIgnoresContextCanceled(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	client := ApiClient{
		logger: logger,
	}

	client.logResponse(
		&http.Request{Method: "POST", URL: &url.URL{Path: "/my-page"}},
		&http.Response{StatusCode: 200},
		context.Canceled,
	)

	logLines := strings.Split(strings.Trim(buf.String(), "\n"), "\n")
	assert.Equal(t, 1, len(logLines))

	logEntries := make([]LogEntry, len(logLines))
	for i, line := range logLines {
		err := json.Unmarshal([]byte(line), &logEntries[i])
		assert.Nil(t, err)
	}

	assert.Equal(t, "INFO", logEntries[0].Level)
	assert.Equal(t, "method: POST, url: /my-page, response: 200", logEntries[0].Message)
}
