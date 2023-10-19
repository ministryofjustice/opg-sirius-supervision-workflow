package urlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSort_GetAriaSort(t *testing.T) {
	assert.Equal(t, "none", Sort{}.GetAriaSort("test"))
	assert.Equal(t, "none", Sort{Descending: true}.GetAriaSort("test"))
	assert.Equal(t, "none", Sort{OrderBy: "foo", Descending: true}.GetAriaSort("test"))
	assert.Equal(t, "ascending", Sort{OrderBy: "test"}.GetAriaSort("test"))
	assert.Equal(t, "descending", Sort{OrderBy: "test", Descending: true}.GetAriaSort("test"))
}

func TestSort_GetDirection(t *testing.T) {
	assert.Equal(t, "asc", Sort{Descending: false}.GetDirection())
	assert.Equal(t, "desc", Sort{Descending: true}.GetDirection())
}

func TestSort_ToURL(t *testing.T) {
	assert.Equal(t, "", Sort{}.ToURL())
	assert.Equal(t, "order-by=test&sort=asc", Sort{OrderBy: "test", Descending: false}.ToURL())
	assert.Equal(t, "order-by=test&sort=desc", Sort{OrderBy: "test", Descending: true}.ToURL())
}
