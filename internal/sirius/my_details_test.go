package sirius

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type request struct {
	method, path string
	cookies      []*http.Cookie
	headers      http.Header
}

func TestMyDetails(t *testing.T) {
	assert := assert.New(t)
	requests := make(chan request, 1)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- request{
			method:  r.Method,
			path:    r.URL.Path,
			cookies: r.Cookies(),
			headers: r.Header,
		}

		io.WriteString(w, `{"id":47,"name":"system","phoneNumber":"03004560300","teams":[{"id":10,"name":"Allocations - (Supervision)","phoneNumber":"0123456789","teams":[],"displayName":"Allocations - (Supervision)","deleted":false,"tasks":[],"email":"allocations.team@opgtest.com"}],"displayName":"system admin","deleted":false,"tasks":[],"email":"system.admin@opgtest.com","firstname":"system","surname":"admin","roles":["OPG User","System Admin"],"locked":false,"suspended":false}`)
	}))
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	cookies := []*http.Cookie{
		{Name: "XSRF-TOKEN", Value: "abcde"},
		{Name: "Other", Value: "other"},
	}

	myDetails, err := client.SiriusUserDetails(getContext(cookies))
	assert.Nil(err)

	assert.Equal(myDetails, UserDetails{
		ID:          47,
		Name:        "system",
		PhoneNumber: "03004560300",
		Teams: []MyDetailsTeam{
			{DisplayName: "Allocations - (Supervision)"},
		},
		DisplayName: "system admin",
		Deleted:     false,
		Email:       "system.admin@opgtest.com",
		Firstname:   "system",
		Surname:     "admin",
		Roles:       []string{"OPG User", "System Admin"},
		Locked:      false,
		Suspended:   false,
	})

	select {
	case r := <-requests:
		assert.Equal(http.MethodGet, r.method)
		assert.Equal("/api/v1/users/current", r.path)
		assert.Equal(cookies, r.cookies)
		assert.Equal("1", r.headers.Get("OPG-Bypass-Membrane"))
		assert.Equal("abcde", r.headers.Get("X-XSRF-TOKEN"))

	case <-time.After(time.Millisecond):
		assert.Fail("request did not happen in time")
	}
}

func TestMyDetailsUnauthorized(t *testing.T) {
	assert := assert.New(t)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusUnauthorized)
	}))
	defer s.Close()

	cookies := []*http.Cookie{
		{Name: "XSRF-TOKEN", Value: "abcde"},
		{Name: "Other", Value: "other"},
	}

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.SiriusUserDetails(getContext(cookies))
	assert.Equal(ErrUnauthorized, err)
}
