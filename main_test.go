package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	assert := assert.New(t)
	port := "9000"

	server := newServer(port)

	go func() {
		server.ListenAndServe()
	}()
	defer server.Shutdown(context.Background())

	resp, err := http.Get("http://localhost:" + port)
	assert.Nil(err)
	assert.Equal(200, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal("Hello Kate and Nick, the beautiful webpage is running", string(data))
}
