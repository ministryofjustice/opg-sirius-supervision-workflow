package urlbuilder

type Filter struct {
	Name                  string
	SelectedValues        []string
	ClearBetweenTeamViews bool
}

func CreateFilter(name string, selectedValues interface{}, clearBetweenTeamViews ...bool) Filter {
	filter := Filter{
		Name: name,
	}

	switch v := selectedValues.(type) {
	case string:
		if v != "" {
			filter.SelectedValues = []string{v}
		}
	case []string:
		filter.SelectedValues = v
	}

	if len(clearBetweenTeamViews) > 0 {
		filter.ClearBetweenTeamViews = clearBetweenTeamViews[0]
	}

	return filter
}
