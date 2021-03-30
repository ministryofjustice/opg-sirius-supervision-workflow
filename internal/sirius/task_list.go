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
	ClientDetailsCaseRecNumber string `json:"caseRecNumber"`
	ClientDetailsFirstName     string `json:"firstname"`
	ClientDetailsId            int    `json:"id"`
	//ClientDetailsMiddlenames          string                     `json:"middlenames"`
	//ClientDetailsSalutation           string                     `json:"salutation"`
	ClientDetailsSupervisionCaseOwner SupervisionCaseOwnerDetail `json:"supervisionCaseOwner"`
	ClientDetailsSurname              string                     `json:"surname"`
	//ClientDetailsUId                  string                     `json:"uId"`
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
	AssigneeDetailsDisplayName string `json:"displayName"`
	//AssigneeDetailsId          int    `json:"id"`
}

type ApiTask struct {
	ApiTaskAssignee  AssigneeDetails    `json:"assignee"`
	ApiTaskCaseItems []CaseItemsDetails `json:"caseItems"`
	// Clients []string `json:"clients"`
	// Description string             `json:"description"`
	ApiTaskDueDate string `json:"dueDate"`
	// ApiTaskId   int                `json:"id"`
	ApiTaskType string `json:"name"`
	// Persons     []string           `json:"persons"`
	// Status      string             `json:"status"`

}

type TaskList struct {
	AllTaskList []ApiTask `json:"tasks"`
}

func (c *Client) GetTaskList(ctx Context, id int) (TaskList, error) {
	var v TaskList
	//Team Id not who is logged in
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/assignees/%d/tasks", id), nil)
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

	allTaskList := v

	return allTaskList, err
}
