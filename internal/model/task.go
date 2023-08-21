package model

import (
	"time"
)

type Task struct {
	Assignee      Assignee `json:"assignee"`
	Orders        []Order  `json:"caseItems"`
	Clients       []Client `json:"clients"`
	Deputies      []Deputy `json:"deputies"`
	DueDate       string   `json:"dueDate"`
	Id            int      `json:"id"`
	Type          string   `json:"type"`
	Name          string   `json:"name"`
	CaseOwnerTask bool     `json:"caseOwnerTask"`
	IsPriority    bool     `json:"isPriority"`
}

type DueDateStatus struct {
	Name   string
	Colour string
}

func (t Task) GetDueDateStatus(now ...time.Time) DueDateStatus {
	removeTime := func(t time.Time) time.Time {
		y, m, d := t.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	}

	today := removeTime(time.Now())
	if len(now) > 0 {
		today = removeTime(now[0])
	}

	dueDate, _ := time.Parse("02/01/2006", t.DueDate)
	dueDate = removeTime(dueDate)

	daysUntilNextWeek := int((7 + (time.Monday - today.Weekday())) % 7)
	if daysUntilNextWeek == 0 {
		daysUntilNextWeek = 7
	}
	startOfNextWeek := today.AddDate(0, 0, daysUntilNextWeek)
	endOfNextWeek := startOfNextWeek.AddDate(0, 0, 6)

	if dueDate.Before(today) {
		return DueDateStatus{Name: "Overdue", Colour: "red"}
	} else if dueDate.Equal(today) {
		return DueDateStatus{Name: "Due Today", Colour: "red"}
	} else if dueDate.Sub(today).Hours() == 24 && dueDate.Before(startOfNextWeek) {
		return DueDateStatus{Name: "Due Tomorrow", Colour: "orange"}
	} else if dueDate.Before(startOfNextWeek) {
		return DueDateStatus{Name: "Due This Week", Colour: "orange"}
	} else if dueDate.Before(endOfNextWeek) || dueDate.Equal(endOfNextWeek) {
		return DueDateStatus{Name: "Due Next Week", Colour: "green"}
	}

	return DueDateStatus{}
}

func (t Task) GetClient() Client {
	if len(t.Orders) != 0 {
		return t.Orders[0].Client
	}
	return t.Clients[0]
}

func (t Task) GetDeputy() Deputy {
	return t.Deputies[0]
}

func (t Task) GetName(taskTypes []TaskType) string {
	for _, taskType := range taskTypes {
		if t.Type == taskType.Handle {
			return taskType.Incomplete
		}
	}
	return t.Name
}

func (t Task) GetAssignee() Assignee {
	if t.Assignee.Name == "Unassigned" {
		if len(t.Clients) != 0 {
			return t.Clients[0].SupervisionCaseOwner
		} else if len(t.Orders) != 0 {
			return t.Orders[0].Client.SupervisionCaseOwner
		}
	}
	return t.Assignee
}
