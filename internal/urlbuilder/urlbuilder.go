package urlbuilder

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"strconv"
	"strings"
)

type UrlBuilder struct {
	Path            string
	SelectedTeam    string
	SelectedPerPage int
	SelectedFilters []Filter
	SelectedSort    Sort
	MyTeamId        string
}

func (ub UrlBuilder) buildUrl(team string, page int, perPage int, filters []Filter, sort Sort, preselectCaseManager bool) string {
	url := ""
	if preselectCaseManager {
		url = fmt.Sprintf("%s?team=%s&page=%d&per-page=%d&preselect", ub.Path, team, page, perPage)
	} else {
		url = fmt.Sprintf("%s?team=%s&page=%d&per-page=%d", ub.Path, team, page, perPage)
	}

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
	return ub.buildUrl(team.Selector, 1, ub.SelectedPerPage, retainedFilters, ub.SelectedSort, CheckIfIsMyTeam(ub.MyTeamId, team.Selector))
}

func CheckIfIsMyTeam(firstId, secondId string) bool {
	return firstId == secondId
}

func (ub UrlBuilder) GetPaginationUrl(page int, perPage ...int) string {
	selectedPerPage := ub.SelectedPerPage
	if len(perPage) > 0 {
		selectedPerPage = perPage[0]
	}
	if strings.HasSuffix(ub.Path, "prefilter") {
		return ub.buildUrl(ub.SelectedTeam, page, selectedPerPage, ub.SelectedFilters, ub.SelectedSort, CheckIfIsMyTeam(ub.MyTeamId, ub.SelectedTeam))
	} else {
		return ub.buildUrl(ub.SelectedTeam, page, selectedPerPage, ub.SelectedFilters, ub.SelectedSort, false)
	}
}

func (ub UrlBuilder) GetSortUrl(orderBy string) string {
	sort := Sort{OrderBy: orderBy}
	if orderBy == ub.SelectedSort.OrderBy {
		sort.Descending = !ub.SelectedSort.Descending
	}
	if strings.HasSuffix(ub.Path, "prefilter") {
		return ub.buildUrl(ub.SelectedTeam, 1, ub.SelectedPerPage, ub.SelectedFilters, sort, CheckIfIsMyTeam(ub.MyTeamId, ub.SelectedTeam))
	} else {
		return ub.buildUrl(ub.SelectedTeam, 1, ub.SelectedPerPage, ub.SelectedFilters, sort, false)
	}
}

func (ub UrlBuilder) GetClearFiltersUrl() string {
	return ub.buildUrl(ub.SelectedTeam, 1, ub.SelectedPerPage, []Filter{}, ub.SelectedSort, false)
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

	return ub.buildUrl(ub.SelectedTeam, 1, ub.SelectedPerPage, retainedFilters, ub.SelectedSort, false), nil
}
