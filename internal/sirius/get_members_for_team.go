package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AssigneesTeam struct {
	Id      int           `json:"id"`
	Members []TeamMembers `json:"members"`
	Name    string        `json:"name"`
}

func (c *Client) GetAssigneesForFilter(ctx Context, loggedInTeamId int, selectedTeam int) (AssigneesTeam, error) {
	var v AssigneesTeam
	teamId := loggedInTeamId

	if selectedTeam != 0 {
		teamId = selectedTeam
	}

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/teams/%d", teamId), nil)

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
