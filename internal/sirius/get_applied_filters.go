package sirius

import (
	"github.com/ministryofjustice/opg-go-common/logging"
)

func (c *Client) GetAppliedFilters(logger *logging.Logger, teamId int, loadTaskTypes []ApiTaskTypes, teamSelection []ReturnedTeamCollection, assigneesForFilter AssigneesTeam) []string {
	var appliedFilters []string

	for _, u := range loadTaskTypes {
		if u.IsSelected {
			appliedFilters = append(appliedFilters, u.Incomplete)
		}
	}

	for _, u := range teamSelection {
		if u.IsTeamSelected && teamId == u.Id {
			appliedFilters = append(appliedFilters, u.Name)
		}
	}

	for _, u := range assigneesForFilter.Members {
		if u.IsSelected {
			appliedFilters = append(appliedFilters, u.TeamMembersDisplayName)
		}
	}
	return appliedFilters
}
