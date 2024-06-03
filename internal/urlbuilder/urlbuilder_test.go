package urlbuilder

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestUrlBuilder_buildUrl(t *testing.T) {
	tests := []struct {
		path    string
		team    string
		page    int
		perPage int
		filters []Filter
		sort    Sort
		want    string
	}{
		{
			path:    "slug",
			team:    "team12",
			page:    11,
			perPage: 25,
			filters: nil,
			want:    "slug?team=team12&page=11&per-page=25",
		},
		{
			path:    "",
			team:    "",
			page:    0,
			perPage: 0,
			filters: []Filter{},
			want:    "?team=&page=0&per-page=0",
		},
		{
			path:    "slug",
			team:    "team12",
			page:    11,
			perPage: 25,
			filters: []Filter{
				{
					Name:           "test",
					SelectedValues: nil,
				},
			},
			want: "slug?team=team12&page=11&per-page=25",
		},
		{
			path:    "slug",
			team:    "team12",
			page:    11,
			perPage: 25,
			filters: []Filter{
				{
					Name:           "test",
					SelectedValues: []string{""},
				},
			},
			want: "slug?team=team12&page=11&per-page=25",
		},
		{
			path:    "slug",
			team:    "team12",
			page:    11,
			perPage: 25,
			filters: []Filter{
				{
					Name:           "test",
					SelectedValues: []string{"val"},
				},
			},
			want: "slug?team=team12&page=11&per-page=25&test=val",
		},
		{
			path:    "slug",
			team:    "team12",
			page:    11,
			perPage: 25,
			filters: []Filter{
				{
					Name:           "test",
					SelectedValues: []string{"val1", "val2"},
				},
			},
			want: "slug?team=team12&page=11&per-page=25&test=val1&test=val2",
		},
		{
			path:    "slug",
			team:    "team12",
			page:    11,
			perPage: 25,
			filters: []Filter{
				{
					Name:           "test",
					SelectedValues: []string{"val1", "val2"},
				},
				{
					Name:           "test2",
					SelectedValues: []string{"val3"},
				},
			},
			want: "slug?team=team12&page=11&per-page=25&test=val1&test=val2&test2=val3",
		},
		{
			path:    "",
			team:    "1",
			page:    2,
			perPage: 15,
			filters: []Filter{},
			sort:    Sort{OrderBy: "name"},
			want:    "?team=1&page=2&per-page=15&order-by=name&sort=asc",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			builder := UrlBuilder{Path: test.path}
			url := builder.buildUrl(test.team, test.page, test.perPage, test.filters, test.sort, false)
			assert.Equal(t, test.want, url)
		})
	}
}

func TestUrlBuilder_GetTeamUrl(t *testing.T) {
	tests := []struct {
		urlBuilder UrlBuilder
		team       string
		want       string
	}{
		{
			urlBuilder: UrlBuilder{Path: "page", SelectedTeam: "lay", SelectedPerPage: 25},
			team:       "lay",
			want:       "page?team=lay&page=1&per-page=25",
		},
		{
			urlBuilder: UrlBuilder{Path: "page", SelectedTeam: "lay", SelectedPerPage: 50},
			team:       "pro",
			want:       "page?team=pro&page=1&per-page=50",
		},
		{
			urlBuilder: UrlBuilder{},
			team:       "pro",
			want:       "?team=pro&page=1&per-page=0",
		},
		{
			urlBuilder: UrlBuilder{SelectedSort: Sort{OrderBy: "name", Descending: true}},
			team:       "pro",
			want:       "?team=pro&page=1&per-page=0&order-by=name&sort=desc",
		},
		{
			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
				{
					Name:                  "cleared1",
					SelectedValues:        []string{"clearedVal1"},
					ClearBetweenTeamViews: true,
				},
				{
					Name:                  "cleared2",
					SelectedValues:        []string{"clearedVal2"},
					ClearBetweenTeamViews: true,
				},
				{
					Name:                  "retained1",
					SelectedValues:        []string{"retainedVal1"},
					ClearBetweenTeamViews: false,
				},
				{
					Name:                  "retained2",
					SelectedValues:        []string{"retainedVal2"},
					ClearBetweenTeamViews: false,
				},
			}},
			team: "pro",
			want: "?team=pro&page=1&per-page=0&retained1=retainedVal1&retained2=retainedVal2",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			team := model.Team{Selector: test.team}
			assert.Equal(t, test.want, test.urlBuilder.GetTeamUrl(team))
		})
	}
}

func TestUrlBuilder_GetPaginationUrl(t *testing.T) {
	tests := []struct {
		urlBuilder UrlBuilder
		page       int
		perPage    []int
		want       string
	}{
		{
			urlBuilder: UrlBuilder{Path: "page", SelectedTeam: "lay", SelectedPerPage: 25},
			page:       2,
			perPage:    []int{25},
			want:       "page?team=lay&page=2&per-page=25",
		},
		{
			urlBuilder: UrlBuilder{SelectedTeam: "lay", SelectedPerPage: 25},
			page:       1,
			perPage:    []int{50},
			want:       "?team=lay&page=1&per-page=50",
		},
		{
			urlBuilder: UrlBuilder{SelectedTeam: "lay", SelectedPerPage: 100},
			page:       2,
			perPage:    nil,
			want:       "?team=lay&page=2&per-page=100",
		},
		{
			urlBuilder: UrlBuilder{SelectedTeam: "lay", SelectedPerPage: 100, SelectedSort: Sort{OrderBy: "name"}},
			page:       2,
			perPage:    nil,
			want:       "?team=lay&page=2&per-page=100&order-by=name&sort=asc",
		},
		{
			urlBuilder: UrlBuilder{SelectedTeam: "lay", SelectedFilters: []Filter{
				{
					Name:                  "retained1",
					SelectedValues:        []string{"val1", "val2"},
					ClearBetweenTeamViews: false,
				},
				{
					Name:                  "retained2",
					SelectedValues:        []string{"val3"},
					ClearBetweenTeamViews: true,
				},
			}},
			page:    2,
			perPage: nil,
			want:    "?team=lay&page=2&per-page=0&retained1=val1&retained1=val2&retained2=val3",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			var result string
			if test.perPage == nil {
				result = test.urlBuilder.GetPaginationUrl(test.page)
			} else {
				result = test.urlBuilder.GetPaginationUrl(test.page, test.perPage[0])
			}
			assert.Equal(t, test.want, result)
		})
	}
}

func TestUrlBuilder_GetClearFiltersUrl(t *testing.T) {
	tests := []struct {
		urlBuilder UrlBuilder
		want       string
	}{
		{
			urlBuilder: UrlBuilder{Path: "page", SelectedTeam: "lay", SelectedPerPage: 50, SelectedSort: Sort{OrderBy: "name"}},
			want:       "page?team=lay&page=1&per-page=50&order-by=name&sort=asc",
		},
		{
			urlBuilder: UrlBuilder{SelectedTeam: "lay", SelectedFilters: []Filter{
				{
					Name:                  "removed1",
					SelectedValues:        []string{"val1"},
					ClearBetweenTeamViews: true,
				},
				{
					Name:                  "removed2",
					SelectedValues:        []string{"val2"},
					ClearBetweenTeamViews: false,
				},
			}},
			want: "?team=lay&page=1&per-page=0",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, test.urlBuilder.GetClearFiltersUrl())
		})
	}
}

func TestUrlBuilder_GetRemoveFilterUrl(t *testing.T) {
	tests := []struct {
		urlBuilder    UrlBuilder
		name          string
		value         interface{}
		want          string
		expectedError error
	}{
		{
			urlBuilder:    UrlBuilder{Path: "page", SelectedTeam: "lay", SelectedPerPage: 50, SelectedSort: Sort{OrderBy: "name"}},
			name:          "non-existent-filter",
			value:         "non-existent-value",
			want:          "page?team=lay&page=1&per-page=50&order-by=name&sort=asc",
			expectedError: nil,
		},
		{
			urlBuilder: UrlBuilder{SelectedTeam: "lay", SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1"},
				},
			}},
			name:          "filter1",
			value:         "non-existent-value",
			want:          "?team=lay&page=1&per-page=0&filter1=val1",
			expectedError: nil,
		},
		{
			urlBuilder: UrlBuilder{SelectedTeam: "lay", SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1"},
				},
			}},
			name:          "filter1",
			value:         "val1",
			want:          "?team=lay&page=1&per-page=0",
			expectedError: nil,
		},
		{
			urlBuilder: UrlBuilder{SelectedTeam: "lay", SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1", "val2"},
				},
			}},
			name:          "filter1",
			value:         "val1",
			want:          "?team=lay&page=1&per-page=0&filter1=val2",
			expectedError: nil,
		},
		{
			urlBuilder: UrlBuilder{SelectedTeam: "lay", SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1", "val2"},
				},
				{
					Name:           "filter2",
					SelectedValues: []string{"val3"},
				},
			}},
			name:          "filter2",
			value:         "val3",
			want:          "?team=lay&page=1&per-page=0&filter1=val1&filter1=val2",
			expectedError: nil,
		},
		{
			urlBuilder: UrlBuilder{SelectedTeam: "lay", SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1", "val2"},
				},
				{
					Name:           "filter2",
					SelectedValues: []string{"23"},
				},
			}},
			name:          "filter2",
			value:         23,
			want:          "?team=lay&page=1&per-page=0&filter1=val1&filter1=val2",
			expectedError: nil,
		},
		{
			urlBuilder: UrlBuilder{SelectedTeam: "lay", SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1", "val2"},
				},
				{
					Name:           "filter2",
					SelectedValues: []string{"23", "45", "66"},
				},
			}},
			name:          "filter2",
			value:         []int{23, 45, 66},
			want:          "",
			expectedError: fmt.Errorf("type []int not accepted"),
		},
		{
			urlBuilder: UrlBuilder{SelectedTeam: "lay", SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1", "val2"},
				},
			}},
			name:          "filter2",
			value:         []string{"val1", "val2", "val3"},
			want:          "",
			expectedError: fmt.Errorf("type []string not accepted"),
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			returnedValue, returnedError := test.urlBuilder.GetRemoveFilterUrl(test.name, test.value)
			assert.Equal(t, test.want, returnedValue)
			assert.Equal(t, test.expectedError, returnedError)
		})
	}
}

func TestUrlBuilder_GetSortUrl(t *testing.T) {
	tests := []struct {
		urlBuilder UrlBuilder
		orderBy    string
		want       string
	}{
		{
			urlBuilder: UrlBuilder{MyTeamId: "9999"},
			orderBy:    "test",
			want:       "?team=&page=1&per-page=0&order-by=test&sort=asc",
		},
		{
			urlBuilder: UrlBuilder{SelectedSort: Sort{OrderBy: "test2", Descending: true}, MyTeamId: "9999"},
			orderBy:    "test",
			want:       "?team=&page=1&per-page=0&order-by=test&sort=asc",
		},
		{
			urlBuilder: UrlBuilder{SelectedSort: Sort{OrderBy: "test"}, MyTeamId: "9999"},
			orderBy:    "test",
			want:       "?team=&page=1&per-page=0&order-by=test&sort=desc",
		},
		{
			urlBuilder: UrlBuilder{SelectedSort: Sort{OrderBy: "test", Descending: true}, MyTeamId: "9999"},
			orderBy:    "test",
			want:       "?team=&page=1&per-page=0&order-by=test&sort=asc",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, test.urlBuilder.GetSortUrl(test.orderBy))
		})
	}
}
