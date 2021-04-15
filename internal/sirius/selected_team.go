package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TeamSelectedMembers struct {
	TeamMembersId   int    `json:"id"`
	TeamMembersName string `json:"name"`
}

type TeamSelected struct {
	Id      int                   `json:"id"`
	Members []TeamSelectedMembers `json:"members"`
	Name    string                `json:"name"`
}

func (c *Client) GetTeamSelected(ctx Context, teamSelection []TeamCollection) (TeamSelected, error) {
	var v TeamSelected
	selectedTeamId := teamSelection[0].UserSelectedTeam

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/teams/%d", selectedTeamId), nil)

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

	return v, err
}
