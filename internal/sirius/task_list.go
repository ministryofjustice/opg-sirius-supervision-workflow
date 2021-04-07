package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SupervisionCaseOwnerDetail struct {
	SupervisionCaseOwnerName string `json:"displayName"`
	//SupervisionCaseOwnerId   int    `json:"id"`
}

type ClientDetails struct {
	ClientCaseRecNumber string `json:"caseRecNumber"`
	ClientFirstName     string `json:"firstname"`
	ClientId            int    `json:"id"`
	//ClientMiddlenames          string                     `json:"middlenames"`
	//ClientSalutation           string                     `json:"salutation"`
	ClientSupervisionCaseOwner SupervisionCaseOwnerDetail `json:"supervisionCaseOwner"`
	ClientSurname              string                     `json:"surname"`
	//ClientUId                  string                     `json:"uId"`
}

type CaseItemsDetails struct {
	//CaseItemCaseRecNumber string        `json:"caseRecNumber"`
	//CaseItemSubtype       string        `json:"caseSubtype"`
	//CaseItemType          string        `json:"caseType"`
	CaseItemClient ClientDetails `json:"client"`
	//CaseItemId            int           `json:"id"`
	//CaseItemUId           string        `json:"uId"`
}

type AssigneeDetails struct {
	AssigneeDisplayName string `json:"displayName"`
	//AssigneeId          int    `json:"id"`
}

type ApiTask struct {
	ApiTaskAssignee  AssigneeDetails    `json:"assignee"`
	ApiTaskCaseItems []CaseItemsDetails `json:"caseItems"`
	// Clients []string `json:"clients"`
	// Description string             `json:"description"`
	ApiTaskDueDate string `json:"dueDate"`
	ApiTaskId      int    `json:"id"`
	ApiTaskType    string `json:"name"`
	// Persons     []string           `json:"persons"`
	// Status      string             `json:"status"`
}

type TaskList struct {
	WholeTaskList []ApiTask `json:"tasks"`
}

func (c *Client) GetTaskList(ctx Context, search int) (TaskList, error) {
	var v TaskList

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/assignees/team/tasks?limit=10&page=%d&sort=dueDate:asc", search), nil)
	if err != nil {
		return v, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, err
	}

	TaskList := v

	return TaskList, err
}
