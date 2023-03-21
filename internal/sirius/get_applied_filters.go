package sirius

func GetAppliedFilters(selectedTeam ReturnedTeamCollection, selectedAssignees []string, selectedUnassigned string, taskTypes []ApiTaskTypes) []string {
	var appliedFilters []string

	for _, u := range taskTypes {
		if u.IsSelected {
			appliedFilters = append(appliedFilters, u.Incomplete)
		}
	}

	if selectedTeam.Selector == selectedUnassigned {
		appliedFilters = append(appliedFilters, selectedTeam.Name)
	}

	for _, u := range selectedTeam.GetAssigneesForFilter() {
		if u.IsSelected(selectedAssignees) {
			appliedFilters = append(appliedFilters, u.Name)
		}
	}

	return appliedFilters
}
