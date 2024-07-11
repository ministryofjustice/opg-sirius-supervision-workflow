package model

import (
	"fmt"
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

func (t Team) GetUnassignedCount(selectedAssignees []AssigneeAndCount, urlPath string) string {
	fmt.Println("in unassigned count")
	for _, a := range selectedAssignees {
		fmt.Println(t.Id)
		if t.Id == a.AssigneeId {
			stringValue := strconv.Itoa(a.Count)
			return "(" + stringValue + ")"
		}
		if a.AssigneeId == 0 {
			stringValue := strconv.Itoa(a.Count)
			return "(" + stringValue + ")"
		}
	}
	//calculate unassigned count for a team of teams
	if t.IsFullLayTeam() || t.IsProDeputyTeam() {
		total := t.GetMultiTeamUnassignedCount(selectedAssignees)
		return "(" + strconv.Itoa(total) + ")"
	}
	return "(0)"
}

func (t Team) GetMultiTeamUnassignedCount(selectedAssignees []AssigneeAndCount) int {
	var total int
	for _, a := range t.Teams {
		total += a.GetCountForATeam(selectedAssignees, a.Id)
	}
	return total
}

func (t Team) GetCountForATeam(selectedAssignees []AssigneeAndCount, teamId int) int {
	for _, a := range selectedAssignees {
		if teamId == a.AssigneeId {
			return a.Count
		}
	}
	return 0
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

func (t Team) IsProDeputyTeam() bool {
	return t.Selector == "pro-team" && t.Id == 0
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

func (t Team) IsLayDeputyTeam() bool {
	return t.Name == "Lay Deputy Team"
}
