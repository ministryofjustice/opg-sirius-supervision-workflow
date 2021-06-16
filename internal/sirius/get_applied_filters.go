package sirius

func (c *Client) GetAppliedFilters(ctx Context, teamId int, loadTaskTypes []ApiTaskTypes, teamSelection []ReturnedTeamCollection, assigneesForFilter AssigneesTeam) []string {
	var appliedFilters []string

	for _, u := range loadTaskTypes {
		if u.IsSelected == true {
			appliedFilters = append(appliedFilters, u.Incomplete)
		}
	}

	for _, u := range teamSelection {
		if u.IsTeamSelected == true && teamId == u.Id {
			appliedFilters = append(appliedFilters, u.Name)
		}
	}

	for _, u := range assigneesForFilter.Members {
		if u.IsSelected == true {
			appliedFilters = append(appliedFilters, u.TeamMembersDisplayName)
		}
	}

	return appliedFilters
}
