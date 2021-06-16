package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type AssigneesTeam struct {
	Id      int                   `json:"id"`
	Members []AssigneeTeamMembers `json:"members"`
	Name    string                `json:"name"`
}
type AssigneeTeamMembers struct {
	TeamMembersId          int    `json:"id"`
	TeamMembersName        string `json:"name"`
	TeamMembersDisplayName string `json:"displayName"`
	IsSelected             bool
}

func (c *Client) GetAssigneesForFilter(ctx Context, teamId int, assigneeSelected []string, appliedFilters []string) (AssigneesTeam, []string, error) {
	var v AssigneesTeam

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/teams/%d", teamId), nil)

	if err != nil {
		return v, []string{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, []string{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, []string{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, []string{}, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, []string{}, err
	}

	assigneeList := make([]AssigneeTeamMembers, len(v.Members))

	for i, u := range v.Members {
		assigneeList[i] = AssigneeTeamMembers{
			TeamMembersId:          u.TeamMembersId,
			TeamMembersName:        u.TeamMembersName,
			TeamMembersDisplayName: u.TeamMembersDisplayName,
			IsSelected:             IsAssigneeSelected(u.TeamMembersId, assigneeSelected),
		}
	}

	for _, u := range assigneeList {
		if u.IsSelected == true {
			appliedFilters = append(appliedFilters, u.TeamMembersDisplayName)
		}
	}

	v.Members = assigneeList

	return v, appliedFilters, err
}

func IsAssigneeSelected(TeamMembersId int, assigneeSelected []string) bool {
	for _, q := range assigneeSelected {
		assigneeSelectedAsAString, _ := strconv.Atoi(q)
		if TeamMembersId == assigneeSelectedAsAString {
			return true
		}
	}
	return false
}
