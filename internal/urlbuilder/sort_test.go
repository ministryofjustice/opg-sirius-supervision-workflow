package urlbuilder

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"strconv"
	"testing"
)

func TestCreateSortFromURL(t *testing.T) {
	tests := []struct {
		urlValues        url.Values
		validSortOptions []string
		want             Sort
	}{
		{
			urlValues:        url.Values{},
			validSortOptions: nil,
			want:             Sort{},
		},
		{
			urlValues:        url.Values{"order-by": {"field"}},
			validSortOptions: nil,
			want:             Sort{},
		},
		{
			urlValues:        url.Values{"order-by": {"field"}},
			validSortOptions: []string{"field"},
			want:             Sort{OrderBy: "field"},
		},
		{
			urlValues:        url.Values{"order-by": {"field"}, "sort": {"invalid"}},
			validSortOptions: []string{"field"},
			want:             Sort{OrderBy: "field"},
		},
		{
			urlValues:        url.Values{"order-by": {"field"}, "sort": {"desc"}},
			validSortOptions: []string{"test", "field"},
			want:             Sort{OrderBy: "field", Descending: true},
		},
		{
			urlValues:        url.Values{"order-by": {"field"}, "sort": {"desc"}},
			validSortOptions: []string{"test"},
			want:             Sort{OrderBy: "test"},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, CreateSortFromURL(test.urlValues, test.validSortOptions))
		})
	}
}

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
