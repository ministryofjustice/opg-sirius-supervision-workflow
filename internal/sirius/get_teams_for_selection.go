package sirius

import (
	"encoding/json"
	"net/http"
)

type TeamMembers struct {
	TeamMembersId          int    `json:"id"`
	TeamMembersName        string `json:"name"`
	TeamMembersDisplayName string `json:"displayName"`
}

type TeamCollection struct {
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

type ReturnedTeamCollection struct {
	Id               int
	Members          []TeamMembers
	Name             string
	UserSelectedTeam int
	SelectedTeamId   int
	Type             string
	TypeLabel        string
}

type TeamStoredData struct {
	TeamId       int
	SelectedTeam int
}

func (c *Client) GetTeamsForSelection(ctx Context, teamId int) ([]ReturnedTeamCollection, error) {
	var v []TeamCollection
	var q []ReturnedTeamCollection
	var k TeamStoredData

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams", nil)
	if err != nil {
		return q, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return q, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return q, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return q, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return q, err
	}

	k.TeamId = teamId

	teams := make([]ReturnedTeamCollection, len(v))
	for i, t := range v {
		teams[i] = ReturnedTeamCollection{
			Id:   t.ID,
			Name: t.DisplayName,
			Type: "",
		}

		for _, m := range t.Members {
			teams[i].Members = append(teams[i].Members, TeamMembers{
				TeamMembersId:   m.ID,
				TeamMembersName: m.Email,
			})
		}

		for i := range teams {
			teams[i].UserSelectedTeam = k.TeamId
			teams[i].SelectedTeamId = k.SelectedTeam

		}
		if t.TeamType != nil {
			teams[i].Type = t.TeamType.Handle
			teams[i].TypeLabel = t.TeamType.Label
		}
	}

	teams = filterOutNonLayTeams(teams)

	return teams, err
}

func filterOutNonLayTeams(v []ReturnedTeamCollection) []ReturnedTeamCollection {
	var filteredTeams []ReturnedTeamCollection
	for _, s := range v {
		if len(s.Type) != 0 {
			filteredTeams = append(filteredTeams, s)
		}
	}
	return filteredTeams
}
