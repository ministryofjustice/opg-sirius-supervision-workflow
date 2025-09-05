package server

import (
	"context"
	"errors"
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
)

func TestRedirect_Error(t *testing.T) {
	assert.Equal(t, "redirect to ", Redirect{Path: ""}.Error())
	assert.Equal(t, "redirect to test-url", Redirect{Path: "test-url"}.Error())
}

func TestRedirect_To(t *testing.T) {
	assert.Equal(t, "", Redirect{Path: ""}.To())
	assert.Equal(t, "test-url", Redirect{Path: "test-url"}.To())
}

func TestStatusError_Code(t *testing.T) {
	assert.Equal(t, 0, StatusError(0).Code())
	assert.Equal(t, 200, StatusError(200).Code())
}

func TestStatusError_Error(t *testing.T) {
	assert.Equal(t, "0 ", StatusError(0).Error())
	assert.Equal(t, "200 OK", StatusError(200).Error())
	assert.Equal(t, "999 ", StatusError(999).Error())
}

type mockNext struct {
	app    WorkflowVars
	w      http.ResponseWriter
	r      *http.Request
	Err    error
	Called int
}

func (m *mockNext) GetHandler() Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		m.app = app
		m.w = w
		m.r = r
		m.Called = m.Called + 1
		return m.Err
	}
}

func recordToMap(rec slog.Record) map[string]interface{} {
	result := make(map[string]interface{})
	rec.Attrs(func(a slog.Attr) bool {
		result[a.Key] = a.Value.Any()
		return true
	})
	return result
}

type TestHandler struct {
	mu      sync.Mutex
	records []slog.Record
}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (h *TestHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *TestHandler) Handle(_ context.Context, rec slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.records = append(h.records, rec.Clone())
	return nil
}

func (h *TestHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *TestHandler) WithGroup(_ string) slog.Handler      { return h }

func (h *TestHandler) Records() []slog.Record {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.records
}

func Test_wrapHandler_error_creating_WorkflowVars(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "test-url", nil)

	mockClient := mockApiClient{error: errors.New("some API error")}

	logHandler := NewTestHandler()
	logger := slog.New(logHandler)

	errorTemplate := &mockTemplate{}
	envVars := EnvironmentVars{}
	cookieStorage := sessions.CookieStore{}
	nextHandlerFunc := wrapHandler(mockClient, logger, errorTemplate, envVars, cookieStorage)
	next := mockNext{}
	httpHandler := nextHandlerFunc(next.GetHandler())
	httpHandler.ServeHTTP(w, r)

	records := logHandler.Records()

	assert.Equal(t, 0, next.Called)
	assert.Len(t, records, 2)
	assert.Equal(t, "Application Request", records[0].Message)
	assert.Equal(t, slog.LevelInfo, records[0].Level)

	assert.Equal(t, "Error handler", records[1].Message)
	assert.Equal(t, slog.LevelError, records[1].Level)
	attrs := recordToMap(records[1])
	assert.Equal(t, "some API error", attrs["error"].(error).Error())

	assert.Equal(t, 1, errorTemplate.count)
	assert.Equal(t, ErrorVars{Code: 500, Error: "some API error"}, errorTemplate.lastVars)
	assert.Equal(t, 500, w.Result().StatusCode)
}

func Test_wrapHandler_status_error_handling(t *testing.T) {
	tests := []struct {
		error     error
		wantCode  int
		wantError string
	}{
		{error: StatusError(400), wantCode: 400, wantError: "400 Bad Request"},
		{error: StatusError(401), wantCode: 401, wantError: "401 Unauthorized"},
		{error: StatusError(403), wantCode: 403, wantError: "403 Forbidden"},
		{error: StatusError(404), wantCode: 404, wantError: "404 Not Found"},
		{error: StatusError(500), wantCode: 500, wantError: "500 Internal Server Error"},
		{error: sirius.StatusError{Code: 400}, wantCode: 400, wantError: "  returned 400"},
		{error: sirius.StatusError{Code: 401}, wantCode: 401, wantError: "  returned 401"},
		{error: sirius.StatusError{Code: 403}, wantCode: 403, wantError: "  returned 403"},
		{error: sirius.StatusError{Code: 404}, wantCode: 404, wantError: "  returned 404"},
		{error: sirius.StatusError{Code: 500}, wantCode: 500, wantError: "  returned 500"},
	}

	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "test-url", nil)

			mockClient := mockApiClient{
				CurrentUserDetails: mockUserDetailsData,
				Teams:              mockTeamsData,
			}

			logHandler := NewTestHandler()
			logger := slog.New(logHandler)

			errorTemplate := &mockTemplate{error: errors.New("some template error")}
			envVars := EnvironmentVars{}
			cookieStorage := sessions.CookieStore{}

			nextHandlerFunc := wrapHandler(mockClient, logger, errorTemplate, envVars, cookieStorage)
			next := mockNext{Err: test.error}
			httpHandler := nextHandlerFunc(next.GetHandler())
			httpHandler.ServeHTTP(w, r)

			records := logHandler.Records()

			assert.Equal(t, 1, next.Called)
			assert.Equal(t, w, next.w)
			assert.Equal(t, r, next.r)

			assert.Len(t, records, 3)
			assert.Equal(t, "Application Request", records[0].Message)
			assert.Equal(t, slog.LevelInfo, records[0].Level)

			assert.Equal(t, "Error handler", records[1].Message)
			assert.Equal(t, slog.LevelError, records[1].Level)
			assert.Equal(t, test.wantError, recordToMap(records[1])["error"].(error).Error())

			assert.Equal(t, "Failed to render error template", records[2].Message)
			assert.Equal(t, "some template error", recordToMap(records[2])["error"].(error).Error())

			assert.Equal(t, 1, errorTemplate.count)
			assert.IsType(t, ErrorVars{}, errorTemplate.lastVars)
			assert.Equal(t, test.wantCode, errorTemplate.lastVars.(ErrorVars).Code)
			assert.Equal(t, test.wantError, errorTemplate.lastVars.(ErrorVars).Error)
			assert.Equal(t, test.wantCode, w.Result().StatusCode)
		})
	}
}

//func Test_wrapHandler_redirects_if_unauthorized(t *testing.T) {
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest(http.MethodGet, "test-url", nil)
//
//	mockClient := mockApiClient{error: sirius.ErrUnauthorized}
//
//	logHandler := NewTestHandler()
//	logger := slog.New(logHandler)
//
//	errorTemplate := &mockTemplate{}
//	envVars := EnvironmentVars{SiriusURL: "sirius-url"}
//	cookieStorage := sessions.CookieStore{}
//
//	nextHandlerFunc := wrapHandler(mockClient, logger, errorTemplate, envVars, cookieStorage)
//	next := mockNext{}
//	httpHandler := nextHandlerFunc(next.GetHandler())
//	httpHandler.ServeHTTP(w, r)
//
//	records := logHandler.Records()
//
//	assert.Equal(t, 0, next.Called)
//	assert.Len(t, records, 1)
//	assert.Equal(t, "Application Request", records[0].Message)
//	assert.Equal(t, "GET", recordToMap(records[0])["method"])
//	assert.Equal(t, "test-url", recordToMap(records[0])["uri"])
//	assert.Equal(t, 0, errorTemplate.count)
//	assert.Equal(t, 302, w.Result().StatusCode)
//
//	location, err := w.Result().Location()
//	assert.Nil(t, err)
//	assert.Equal(t, "sirius-url/auth", location.String())
//}
//
//func Test_wrapHandler_follows_local_redirect(t *testing.T) {
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest(http.MethodGet, "test-url", nil)
//
//	mockClient := mockApiClient{
//		CurrentUserDetails: mockUserDetailsData,
//		Teams:              mockTeamsData,
//	}
//
//	logHandler := NewTestHandler()
//	logger := slog.New(logHandler)
//
//	errorTemplate := &mockTemplate{}
//	envVars := EnvironmentVars{Prefix: "/workflow-prefix"}
//	cookieStorage := sessions.CookieStore{}
//
//	nextHandlerFunc := wrapHandler(mockClient, logger, errorTemplate, envVars, cookieStorage)
//	next := mockNext{Err: Redirect{Path: "redirect-to-here"}}
//	httpHandler := nextHandlerFunc(next.GetHandler())
//	httpHandler.ServeHTTP(w, r)
//
//	records := logHandler.Records()
//
//	assert.Equal(t, 1, next.Called)
//	assert.Equal(t, w, next.w)
//	assert.Equal(t, r, next.r)
//	assert.Len(t, records, 1)
//	assert.Equal(t, "Application Request", records[0].Message)
//	assert.Equal(t, "GET", recordToMap(records[0])["method"])
//	assert.Equal(t, "test-url", recordToMap(records[0])["uri"])
//	assert.Equal(t, 0, errorTemplate.count)
//	assert.Equal(t, 302, w.Result().StatusCode)
//
//	location, err := w.Result().Location()
//	assert.Nil(t, err)
//	assert.Equal(t, "/workflow-prefix/redirect-to-here", location.String())
//}
//
//func Test_wrapHandler_leaves_canceled_context_early(t *testing.T) {
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest(http.MethodGet, "test-url", nil)
//
//	mockClient := mockApiClient{error: context.Canceled}
//
//	logHandler := NewTestHandler()
//	logger := slog.New(logHandler)
//
//	errorTemplate := &mockTemplate{}
//	envVars := EnvironmentVars{SiriusURL: "sirius-url"}
//	cookieStorage := sessions.CookieStore{}
//
//	nextHandlerFunc := wrapHandler(mockClient, logger, errorTemplate, envVars, cookieStorage)
//	next := mockNext{}
//	httpHandler := nextHandlerFunc(next.GetHandler())
//	httpHandler.ServeHTTP(w, r)
//
//	records := logHandler.Records()
//
//	assert.Equal(t, 0, next.Called)
//	assert.Len(t, records, 1)
//	assert.Equal(t, "Application Request", records[0].Message)
//	assert.Equal(t, 0, errorTemplate.count)
//	assert.Equal(t, 499, w.Result().StatusCode)
//}
