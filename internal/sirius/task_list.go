package sirius

import (
	"encoding/json"
	"net/http"
)

type SupervisionCaseOwnerDetail struct {
	DisplayName            string `json:"displayName"`
	SupervisionCaseOwnerId int    `json:"id"`
}

type ClientDetails struct {
	CaseRecNumber        string                     `json:"caseRecNumber"`
	TaskFirstname        string                     `json:"firstname"`
	ClientId             int                        `json:"id"`
	ClientMiddlenames    string                     `json:"middlenames"`
	ClientSalutation     string                     `json:"salutation"`
	SupervisionCaseOwner SupervisionCaseOwnerDetail `json:"supervisionCaseOwner"`
	TaskSurname          string                     `json:"surname"`
	ClientUId            string                     `json:"uId"`
}

type CaseItemsDetails struct {
	CaseRecNumber string        `json:"caseRecNumber"`
	CaseSubtype   string        `json:"caseSubtype"`
	CaseType      string        `json:"caseType"`
	Client        ClientDetails `json:"client"`
	CaseItemsId   int           `json:"id"`
	CaseItemsUId  string        `json:"uId"`
}

type AssigneeDetails struct {
	DisplayName string `json:"displayName"`
	AssigneeId  int    `json:"id"`
}

type ApiTask struct {
	Assignee    AssigneeDetails    `json:"assignee"`
	CaseItems   []CaseItemsDetails `json:"caseItems"`
	Clients     []string           `json:"clients"`
	Description string             `json:"description"`
	DueDate     string             `json:"dueDate"`
	ApiTaskId   int                `json:"id"`
	Name        string             `json:"name"`
	Persons     []string           `json:"persons"`
	Status      string             `json:"status"`
}

type TaskList struct {
	AllTaskList []ApiTask `json:"tasks"`
}

func (c *Client) GetTaskList(ctx Context) ([]ApiTask, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/assignees/65/tasks", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v TaskList
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	allTaskList := v.AllTaskList

	return allTaskList, err
}
