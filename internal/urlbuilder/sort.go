package urlbuilder

import (
	"net/url"
	"slices"
)

type Sort struct {
	OrderBy    string
	Descending bool
}

func CreateSortFromURL(values url.Values, validOptions []string) Sort {
	if len(validOptions) == 0 {
		return Sort{}
	}
	sort := Sort{
		OrderBy:    values.Get("order-by"),
		Descending: values.Get("sort") == "desc",
	}
	if slices.Contains(validOptions, sort.OrderBy) {
		return sort
	}
	return Sort{OrderBy: validOptions[0]}
}

func (s Sort) GetDirection() string {
	if s.Descending {
		return "desc"
	}
	return "asc"
}

func (s Sort) ToURL() string {
	if s.OrderBy == "" {
		return ""
	}
	values := url.Values{}
	values.Set("order-by", s.OrderBy)
	values.Set("sort", s.GetDirection())
	return values.Encode()
}

func (s Sort) GetAriaSort(orderBy string) string {
	if s.OrderBy != orderBy {
		return "none"
	}
	if s.Descending {
		return "descending"
	}
	return "ascending"
}
