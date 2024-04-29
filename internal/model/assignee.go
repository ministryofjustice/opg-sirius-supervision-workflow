package model

import (
	"strconv"
	"strings"
)

type AssigneeAndCount struct {
	AssigneeId int `json:"assignee"`
	Count      int `json:"count"`
}

type Assignee struct {
	Id          int      `json:"id"`
	Name        string   `json:"displayName"`
	Teams       []Team   `json:"teams"`
	PhoneNumber string   `json:"phoneNumber"`
	Deleted     bool     `json:"deleted"`
	Email       string   `json:"email"`
	Firstname   string   `json:"firstname"`
	Surname     string   `json:"surname"`
	Roles       []string `json:"roles"`
	Locked      bool     `json:"locked"`
	Suspended   bool     `json:"suspended"`
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
	if urlPath == "deputy-tasks" {
		return ""
	}
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
