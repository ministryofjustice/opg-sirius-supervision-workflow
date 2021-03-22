package sirius

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPermissionStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.MyPermissions(getContext(nil))
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/permissions",
		Method: http.MethodGet,
	}, err)
}

func TestPermissionSetChecksPermission(t *testing.T) {
	permissions := PermissionSet{
		"user": {
			Permissions: []string{"GET", "PATCH"},
		},
		"team": {
			Permissions: []string{"GET"},
		},
	}

	assert.True(t, permissions.HasPermission("user", "PATCH"))
	assert.True(t, permissions.HasPermission("team", "GET"))
	assert.True(t, permissions.HasPermission("team", "get"))
	assert.False(t, permissions.HasPermission("team", "PATCHs"))
}
