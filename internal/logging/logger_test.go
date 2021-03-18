package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	logger := New(&buf, "hi")

	logger.Print("one", "two", "three")

	var v logEvent
	assert.Nil(json.NewDecoder(&buf).Decode(&v))

	assert.Equal("hi", v.ServiceName)
	assert.WithinDuration(time.Now(), v.Timestamp, time.Second)
	assert.Equal("onetwothree", v.Message)
}

func TestRequest(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	logger := New(&buf, "hi")
	r, _ := http.NewRequest("GET", "/something", nil)

	logger.Request(r, errors.New("what"))

	var v requestEvent
	assert.Nil(json.NewDecoder(&buf).Decode(&v))

	assert.Equal("hi", v.ServiceName)
	assert.WithinDuration(time.Now(), v.Timestamp, time.Second)
	assert.Equal("GET", v.RequestMethod)
	assert.Equal("/something", v.RequestURI)
	assert.Equal("what", v.Message)
	assert.Nil(v.Data)
}

type anExpandedError struct {
	title string
	data  interface{}
}

func (e anExpandedError) Error() string     { return "impl" }
func (e anExpandedError) Title() string     { return e.title }
func (e anExpandedError) Data() interface{} { return e.data }

func TestRequestWithExpandedError(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	logger := New(&buf, "hi")
	r, _ := http.NewRequest("GET", "/something", nil)

	err := anExpandedError{
		title: "a title",
		data:  map[string]interface{}{"hello": "there"},
	}

	logger.Request(r, err)

	var v requestEvent
	assert.Nil(json.NewDecoder(&buf).Decode(&v))

	assert.Equal("hi", v.ServiceName)
	assert.WithinDuration(time.Now(), v.Timestamp, time.Second)
	assert.Equal("GET", v.RequestMethod)
	assert.Equal("/something", v.RequestURI)
	assert.Equal(err.title, v.Message)
	assert.Equal(err.data, v.Data)
}
