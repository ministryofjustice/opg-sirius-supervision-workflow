package sirius

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
)

type TeamMember struct {
	ID   int    `json:"id"`
	Name string `json:"displayName"`
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

type Team struct {
	Id        int
	Members   []TeamMember
	Name      string
	Type      string
	TypeLabel string
	Selector  string
	Teams     []Team
}

func (t Team) GetAssigneesForFilter() []TeamMember {
	assignees := t.Members
	for _, team := range t.Teams {
		assignees = append(assignees, team.Members...)
	}
	ids := map[int]bool{}
	var deduped []TeamMember
	for _, assignee := range assignees {
		if _, value := ids[assignee.ID]; !value {
			ids[assignee.ID] = true
			deduped = append(deduped, assignee)
		}
	}
	sort.Slice(deduped, func(i, j int) bool {
		return deduped[i].Name < deduped[j].Name
	})
	return deduped
}

func (t Team) HasTeam(id int) bool {
	if t.Id == id {
		return true
	}
	for _, t := range t.Teams {
		if t.Id == id {
			return true
		}
	}
	return false
}

func (t Team) IsLay() bool {
	return t.Type == "LAY" || t.Selector == "lay-team"
}

func (t Team) IsPro() bool {
	return t.Type == "PRO" || t.Selector == "pro-team"
}

func (m TeamMember) IsSelected(selectedAssignees []string) bool {
	for _, a := range selectedAssignees {
		id, _ := strconv.Atoi(a)
		if m.ID == id {
			return true
		}
	}
	return false
}

func (c *ApiClient) GetTeamsForSelection(ctx Context) ([]Team, error) {
	var v []TeamCollection
	var q []Team

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

	layTeam := Team{
		Members:  []TeamMember{},
		Name:     "Lay deputy team",
		Selector: "lay-team",
		Teams:    []Team{},
	}

	proTeam := Team{
		Members:  []TeamMember{},
		Name:     "Professional deputy team",
		Selector: "pro-team",
		Teams:    []Team{},
	}

	for _, t := range v {
		if t.TeamType == nil {
			continue
		}

		team := Team{
			Id:        t.ID,
			Name:      t.DisplayName,
			Type:      t.TeamType.Handle,
			TypeLabel: t.TeamType.Label,
			Selector:  strconv.Itoa(t.ID),
			Teams:     []Team{},
		}

		for _, m := range t.Members {
			team.Members = append(team.Members, TeamMember{
				ID:   m.ID,
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

	return q, err
}
