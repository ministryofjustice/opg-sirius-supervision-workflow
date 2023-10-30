package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"net/http"
	"sort"
	"strconv"
)

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

func (c *ApiClient) GetTeamsForSelection(ctx Context, teamTypes []string) ([]model.Team, error) {
	var v []TeamCollection
	var q []model.Team

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams", nil)
	if err != nil {
		return q, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.logErrorRequest(req, err)
		return q, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.logResponse(req, resp, err)
		return q, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logResponse(req, resp, err)
		return q, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		c.logResponse(req, resp, err)
		return q, err
	}

	layTeam := model.Team{
		Members:  []model.Assignee{},
		Name:     "Lay Deputy Team",
		Selector: "lay-team",
		Teams:    []model.Team{},
	}

	proTeam := model.Team{
		Members:  []model.Assignee{},
		Name:     "Professional Deputy Team",
		Selector: "pro-team",
		Teams:    []model.Team{},
	}

	for _, t := range v {
		if t.TeamType == nil {
			continue
		}

		team := model.Team{
			Id:        t.ID,
			Name:      t.DisplayName,
			Type:      t.TeamType.Handle,
			TypeLabel: t.TeamType.Label,
			Selector:  strconv.Itoa(t.ID),
			Teams:     []model.Team{},
		}

		for _, m := range t.Members {
			team.Members = append(team.Members, model.Assignee{
				Id:   m.ID,
				Name: m.DisplayName,
			})
		}

		if team.IsLay() {
			layTeam.Members = append(layTeam.Members, team.Members...)
			layTeam.Teams = append(layTeam.Teams, team)
		} else if team.IsPro() {
			proTeam.Members = append(proTeam.Members, team.Members...)
			proTeam.Teams = append(proTeam.Teams, team)
		}

		q = append(q, team)
	}

	q = append(q, layTeam, proTeam)

	sort.Slice(q, func(i, j int) bool {
		return q[i].Name < q[j].Name
	})

	if len(teamTypes) > 0 {
		q = returnSelectedTeamTypes(q, teamTypes)
	}
	fmt.Println("q")
	fmt.Println(q)
	return q, err
}

func returnSelectedTeamTypes(allTeams []model.Team, teamTypes []string) []model.Team {
	var teamsToReturn []model.Team

	for _, tt := range teamTypes {
		for i, m := range allTeams {
			if m.Type == tt {
				teamsToReturn = append(teamsToReturn, allTeams[i])
				//if len(allTeams[i].Members) > 0 {
				//	teamsToReturn[i].Members = allTeams[i].Members
				//}
			}
		}
	}

	return teamsToReturn
}
