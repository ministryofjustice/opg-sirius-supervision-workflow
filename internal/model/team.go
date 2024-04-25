package model

import (
	"sort"
	"strconv"
	"strings"
)

type Team struct {
	Id        int    `json:"id"`
	Name      string `json:"displayName"`
	Members   []Assignee
	Type      string
	TypeLabel string
	Selector  string
	Teams     []Team
}

func (t Team) GetAssigneesForFilter() []Assignee {
	assignees := t.Members
	for _, team := range t.Teams {
		assignees = append(assignees, team.Members...)
	}
	ids := map[int]bool{}
	var deduped []Assignee
	for _, assignee := range assignees {
		if _, value := ids[assignee.Id]; !value {
			ids[assignee.Id] = true
			deduped = append(deduped, assignee)
		}
	}
	sort.Slice(deduped, func(i, j int) bool {
		return deduped[i].Name < deduped[j].Name
	})
	return deduped
}

func (t Team) GetUnassignedCount(selectedAssignees []AssigneeAndCount) string {
	for _, a := range selectedAssignees {
		if t.Id == a.AssigneeId {
			stringValue := strconv.Itoa(a.Count)
			return "(" + stringValue + ")"
		}
	}
	return "(0)"
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

func (t Team) IsFullLayTeam() bool {
	return t.Selector == "lay-team"
}

func (t Team) IsLay() bool {
	return t.Type == "LAY" || t.IsFullLayTeam()
}

func (t Team) IsPro() bool {
	return t.Type == "PRO" || t.Selector == "pro-team"
}

func (t Team) IsPA() bool {
	return t.Type == "PA"
}

func (t Team) IsHW() bool {
	return t.Type == "HW"
}

func (t Team) IsClosedCases() bool {
	return strings.ToLower(t.Name) == "supervision closed cases"
}

func (t Team) IsLayNewOrdersTeam() bool {
	return t.Name == "Lay Team - New Deputy Orders"
}
