package sirius

import (
	"encoding/json"
	"net/http"
)

type apiTeam struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Members     []struct {
		ID          int    `json:"id"`
		DisplayName string `json:"displayName"`
		Email       string `json:"email"`
	} `json:"members"`
	TeamType *struct {
		Handle string `json:"handle"`
		Label  string `json:"label"`
	} `json:"teamType"`
}

type TeamMember struct {
	ID          int
	DisplayName string
	Email       string
}

type Team struct {
	ID          int
	DisplayName string
	Members     []TeamMember
	Type        string
	TypeLabel   string
	Email       string
	PhoneNumber string
}

func (c *Client) Teams(ctx Context) ([]Team, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v []apiTeam
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	teams := make([]Team, len(v))

	for i, t := range v {
		teams[i] = Team{
			ID:          t.ID,
			DisplayName: t.DisplayName,
			Type:        "",
			TypeLabel:   "LPA",
		}

		for _, m := range t.Members {
			teams[i].Members = append(teams[i].Members, TeamMember{
				DisplayName: m.DisplayName,
				Email:       m.Email,
			})
		}

		if t.TeamType != nil {
			teams[i].Type = t.TeamType.Handle
			teams[i].TypeLabel = "Supervision — " + t.TeamType.Label
		}
	}

	return teams, nil
}
