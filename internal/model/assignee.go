package model

import (
	"strconv"
	"strings"
)

// AssigneeAndCount is decoded directly from JSON (no DTO layer) since it is a
// small, non-recursive struct used as-is by both the API response and the domain model.
type AssigneeAndCount struct {
	AssigneeId int `json:"assignee"`
	Count      int `json:"count"`
}

type Assignee struct {
	Id          int
	Name        string
	Teams       []Team
	PhoneNumber string
	Deleted     bool
	Email       string
	Firstname   string
	Surname     string
	Roles       []string
	Locked      bool
	Suspended   bool
}

func (m Assignee) IsSelected(selectedAssignees []string) bool {
	for _, a := range selectedAssignees {
		id, _ := strconv.Atoi(a)
		if m.Id == id {
			return true
		}
	}
	return false
}

func (m Assignee) GetCountAsString(selectedAssignees []AssigneeAndCount, urlPath string) string {
	for _, a := range selectedAssignees {

		if m.Id == a.AssigneeId {
			stringValue := strconv.Itoa(a.Count)
			return "(" + stringValue + ")"
		}
	}
	return "(0)"
}

func (m Assignee) GetRoles() string {
	return strings.Join(m.Roles, ",")
}

func (m Assignee) IsOnlyCaseManager() bool {
	//allow for 2 roles as one will always be OPG User
	if len(m.Roles) < 3 {
		for _, a := range m.Roles {
			if strings.ToLower(a) == "case manager" {
				return true
			}
		}
	}
	return false
}
