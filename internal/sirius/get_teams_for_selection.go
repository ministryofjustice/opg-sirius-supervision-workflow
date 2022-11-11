package sirius

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-go-common/logging"
	"net/http"
	"strconv"
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
	IsTeamSelected   bool
}

type TeamStoredData struct {
	TeamId       int
	SelectedTeam int
}

func (c *Client) GetTeamsForSelection(ctx Context, logger *logging.Logger, teamId int, assigneeSelected []string) ([]ReturnedTeamCollection, error) {
	var v []TeamCollection
	var q []ReturnedTeamCollection
	var k TeamStoredData

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams", nil)
	c.logRequest(req, err)
	if err != nil {
		return q, err
	}

	resp, err := c.http.Do(req)
	c.logResponse(resp, err)

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
			Id:             t.ID,
			Name:           t.DisplayName,
			Type:           "",
			IsTeamSelected: IsTeamSelected(teamId, assigneeSelected, t.ID),
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

	teams = FilterOutNonLayTeams(teams)

	return teams, err
}

func FilterOutNonLayTeams(v []ReturnedTeamCollection) []ReturnedTeamCollection {
	var filteredTeams []ReturnedTeamCollection
	for _, s := range v {
		if len(s.Type) != 0 {
			filteredTeams = append(filteredTeams, s)
		}
	}
	return filteredTeams
}

func IsTeamSelected(teamId int, assigneeSelected []string, myTeamId int) bool {
	for _, q := range assigneeSelected {
		assigneeSelectedAsAString, _ := strconv.Atoi(q)
		if teamId == assigneeSelectedAsAString && teamId == myTeamId {
			return true
		}
	}
	return false
}
