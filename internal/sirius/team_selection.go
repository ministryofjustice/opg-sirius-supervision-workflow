package sirius

import (
	"encoding/json"
	"net/http"
)

type TeamMembers struct {
	TeamMembersId   int    `json:"id"`
	TeamMembersName string `json:"name"`
}

type TeamCollection struct {
	Id               int           `json:"id"`
	Members          []TeamMembers `json:"members"`
	Name             string        `json:"name"`
	UserSelectedTeam int
	SelectedTeamId   int
}

type TeamStoredData struct {
	TeamId       int
	SelectedTeam int
}

func (c *Client) GetTeamSelection(ctx Context, myDetails UserDetails, selectedTeamName int, selectedTeamMembers TeamSelected) ([]TeamCollection, error) {
	var v []TeamCollection
	var k TeamStoredData

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams", nil)

	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, err
	}

	if selectedTeamName == 0 && k.TeamId == 0 {
		k.TeamId = myDetails.Teams[0].TeamId
	} else {
		k.TeamId = selectedTeamName
	}

	if selectedTeamMembers.selectedTeamToAssignTask == 0 && k.SelectedTeam == 0 {
		k.SelectedTeam = myDetails.Teams[0].TeamId
	} else {
		k.SelectedTeam = selectedTeamMembers.selectedTeamToAssignTask
	}

	for i, _ := range v {
		v[i].UserSelectedTeam = k.TeamId
		v[i].SelectedTeamId = k.SelectedTeam
	}

	return v, err
}
