package urlbuilder

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"strconv"
)

type UrlBuilder struct {
	Path            string
	SelectedTeam    string
	SelectedPerPage int
	SelectedFilters []Filter
	SelectedSort    Sort
}

func (ub UrlBuilder) buildUrl(team string, page int, perPage int, filters []Filter, sort Sort) string {
	url := fmt.Sprintf("%s?team=%s&page=%d&per-page=%d", ub.Path, team, page, perPage)
	for _, filter := range filters {
		for _, value := range filter.SelectedValues {
			if value != "" {
				url += "&" + filter.Name + "=" + value
			}
		}
	}
	if sort.ToURL() != "" {
		url += "&" + sort.ToURL()
	}
	return url
}

func (ub UrlBuilder) GetTeamUrl(team model.Team) string {
	var retainedFilters []Filter
	for _, filter := range ub.SelectedFilters {
		if !filter.ClearBetweenTeamViews {
			retainedFilters = append(retainedFilters, filter)
		}
	}
	return ub.buildUrl(team.Selector, 1, ub.SelectedPerPage, retainedFilters, ub.SelectedSort)
}

func (ub UrlBuilder) GetPaginationUrl(page int, perPage ...int) string {
	selectedPerPage := ub.SelectedPerPage
	if len(perPage) > 0 {
		selectedPerPage = perPage[0]
	}
	return ub.buildUrl(ub.SelectedTeam, page, selectedPerPage, ub.SelectedFilters, ub.SelectedSort)
}

func (ub UrlBuilder) GetSortUrl(orderBy string) string {
	sort := Sort{OrderBy: orderBy}
	if orderBy == ub.SelectedSort.OrderBy {
		sort.Descending = !ub.SelectedSort.Descending
	}
	return ub.buildUrl(ub.SelectedTeam, 1, ub.SelectedPerPage, ub.SelectedFilters, sort)
}

func (ub UrlBuilder) GetClearFiltersUrl() string {
	return ub.buildUrl(ub.SelectedTeam, 1, ub.SelectedPerPage, []Filter{}, ub.SelectedSort)
}

func (ub UrlBuilder) GetRemoveFilterUrl(name string, value interface{}) (string, error) {
	var retainedFilters []Filter
	var retainedValues []string
	var stringValue string

	switch val := value.(type) {
	case string:
		stringValue = val
	case int:
		stringValue = strconv.Itoa(val)
	default:
		return "", fmt.Errorf("type %T not accepted", val)
	}

	for _, filter := range ub.SelectedFilters {
		retainedValues = []string{}
		for _, v := range filter.SelectedValues {
			if name != filter.Name || stringValue != v {
				retainedValues = append(retainedValues, v)
			}
		}
		if len(retainedValues) > 0 {
			filter.SelectedValues = retainedValues
			retainedFilters = append(retainedFilters, filter)
		}
	}

	return ub.buildUrl(ub.SelectedTeam, 1, ub.SelectedPerPage, retainedFilters, ub.SelectedSort), nil
}
