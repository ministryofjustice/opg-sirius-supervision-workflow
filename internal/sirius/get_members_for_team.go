package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TeamSelectedMembers struct {
	TeamMembersId   int    `json:"id"`
	TeamMembersName string `json:"displayName"`
}

type TeamSelected struct {
	Id                       int                   `json:"id"`
	Members                  []TeamSelectedMembers `json:"members"`
	Name                     string                `json:"name"`
	selectedTeamToAssignTask int
}

// {14 [{106 LayTeam2 User1}] Lay Team 2 - (Supervision) 14}

func (c *Client) GetMembersForTeam(ctx Context, loggedInTeamId int, selectedTeamToAssignTask int) (TeamSelected, error) {
	//13 14
	var v TeamSelected

	if selectedTeamToAssignTask == 0 && v.selectedTeamToAssignTask == 0 {
		selectedTeamToAssignTask = loggedInTeamId
	}

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/teams/%d", selectedTeamToAssignTask), nil)

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

	v.selectedTeamToAssignTask = selectedTeamToAssignTask

	return v, err
}
