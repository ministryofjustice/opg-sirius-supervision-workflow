package sirius

import "time"

func GetAppliedFilters(selectedTeam ReturnedTeamCollection, selectedAssignees []string, selectedUnassigned string, taskTypes []TaskTypes, dueDateFrom *time.Time, dueDateTo *time.Time) []string {
	var appliedFilters []string

	for _, u := range taskTypes {
		if u.IsSelected {
			appliedFilters = append(appliedFilters, u.Name)
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

	if dueDateFrom != nil {
		appliedFilters = append(appliedFilters, "Due date from "+dueDateFrom.Format("02/01/2006")+" (inclusive)")
	}

	if dueDateTo != nil {
		appliedFilters = append(appliedFilters, "Due date to "+dueDateTo.Format("02/01/2006")+" (inclusive)")
	}

	return appliedFilters
}
