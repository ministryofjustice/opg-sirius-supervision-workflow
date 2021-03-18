package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Logger struct {
	serviceName string

	mu  sync.Mutex
	out *json.Encoder
}

func New(out io.Writer, serviceName string) *Logger {
	return &Logger{out: json.NewEncoder(out), serviceName: serviceName}
}

type logEvent struct {
	ServiceName string    `json:"service_name"`
	Timestamp   time.Time `json:"timestamp"`
	Message     string    `json:"message"`
}

func (l *Logger) Print(v ...interface{}) {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	_ = l.out.Encode(logEvent{
		ServiceName: l.serviceName,
		Message:     fmt.Sprint(v...),
		Timestamp:   now,
	})
}

func (l *Logger) Fatal(err error) {
	l.Print(err)
	os.Exit(1)
}

type requestEvent struct {
	ServiceName   string      `json:"service_name"`
	Timestamp     time.Time   `json:"timestamp"`
	RequestMethod string      `json:"request_method"`
	RequestURI    string      `json:"request_uri"`
	Message       string      `json:"message"`
	Data          interface{} `json:"data"`
}

type expandedError interface {
	Title() string
	Data() interface{}
}

func (l *Logger) Request(r *http.Request, err error) {
	now := time.Now()

	event := requestEvent{
		ServiceName:   l.serviceName,
		RequestMethod: r.Method,
		RequestURI:    r.URL.String(),
		Message:       err.Error(),
		Timestamp:     now,
	}

	if ee, ok := err.(expandedError); ok {
		event.Message = ee.Title()
		event.Data = ee.Data()
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	_ = l.out.Encode(event)
}
