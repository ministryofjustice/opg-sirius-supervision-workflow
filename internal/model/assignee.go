package model

import (
	"strconv"
	"strings"
)

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
	TaskCount   int
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

//
//func (m Assignee) GetTaskCount(assigneeCount []sirius.AssigneeAndCount) int {
//	for i := 0; i < len(assigneeCount); i++ {
//		if m.Id == assigneeCount[i].AssigneeId {
//			return assigneeCount[i].Count
//		}
//	}
//	return 0
//}

func (m Assignee) GetRoles() string {
	return strings.Join(m.Roles, ",")
}
