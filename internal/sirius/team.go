package sirius

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (c *Client) Team(ctx Context, id int) (Team, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams/"+strconv.Itoa(id), nil)
	if err != nil {
		return Team{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Team{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return Team{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return Team{}, newStatusError(resp)
	}

	var v apiTeam
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return Team{}, err
	}

	team := Team{
		ID:          v.ID,
		DisplayName: v.DisplayName,
		Type:        "",
		Email:       v.Email,
		PhoneNumber: v.PhoneNumber,
	}

	for _, m := range v.Members {
		team.Members = append(team.Members, TeamMember{
			ID:          m.ID,
			DisplayName: m.DisplayName,
			Email:       m.Email,
		})
	}

	if v.TeamType != nil {
		team.Type = v.TeamType.Handle
	}

	return team, nil
}
