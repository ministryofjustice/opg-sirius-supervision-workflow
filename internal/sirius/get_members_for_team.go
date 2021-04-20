package sirius

import (
	"encoding/json"
	"fmt"
	"log"
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

func (c *Client) GetMembersForTeam(ctx Context, myDetails UserDetails, selectedTeamToAssignTask int) (TeamSelected, error) {
	var v TeamSelected
	log.Println("v.selectedTeamToAssignTask")
	log.Println(v.selectedTeamToAssignTask)
	log.Println("selectedTeamToAssignTask")
	log.Println(selectedTeamToAssignTask)
	if selectedTeamToAssignTask == 0 && v.selectedTeamToAssignTask == 0 {
		selectedTeamToAssignTask = myDetails.Teams[0].TeamId
	}

	// log.Println(selectedTeamToAssignTask)
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
