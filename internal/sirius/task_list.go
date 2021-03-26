package sirius

import (
	"encoding/json"
	"net/http"
)

type supervisionCaseOwnerDetail struct {
	DisplayName            string `json:"displayName"`
	SupervisionCaseOwnerId int    `json:"id"`
}

type clientDetails struct {
	CaseRecNumber        string                     `json:"caseRecNumber"`
	TaskFirstname        string                     `json:"firstname"`
	ClientId             int                        `json:"id"`
	ClientMiddlenames    string                     `json:"middlenames"`
	ClientSalutation     string                     `json:"salutation"`
	SupervisionCaseOwner supervisionCaseOwnerDetail `json:"supervisionCaseOwner"`
	TaskSurname          string                     `json:"surname"`
	ClientUId            string                     `json:"uId"`
}

type caseItemsDetails struct {
	CaseRecNumber string        `json:"caseRecNumber"`
	CaseSubtype   string        `json:"caseSubtype"`
	CaseType      string        `json:"caseType"`
	Client        clientDetails `json:"client"`
	CaseItemsId   int           `json:"id"`
	CaseItemsUId  string        `json:"uId"`
}

type AssigneeDetails struct {
	DisplayName string `json:"displayName"`
	AssigneeId  int    `json:"id"`
}

type ApiTask struct {
	Assignee    AssigneeDetails    `json:"assignee"`
	CaseItems   []caseItemsDetails `json:"caseItems"`
	Clients     []string           `json:"clients"`
	Description string             `json:"description"`
	DueDate     string             `json:"dueDate"`
	ApiTaskId   int                `json:"id"`
	Name        string             `json:"name"`
	Persons     []string           `json:"persons"`
	Status      string             `json:"status"`
}

type TaskList struct {
	AllTaskList []ApiTask `json:"tasks"` //look into the type of this map for next time
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

	// io.Copy(os.Stdout, resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v TaskList
	// var v []ApiTask
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	allTaskList := v.AllTaskList

	// allTaskList := make([]MyTaskList, len(v))

	// var taskList []ApiTask

	// for _, u := range allTaskList {
	// 	task := ApiTask{
	// 		// Assignee:    u.Assignee,
	// 		// CaseItems:   u.CaseItems,
	// 		// Clients:     u.Clients,
	// 		// CreatedTime: u.CreatedTime,
	// 		// Description: u.Description,
	// 		// DueDate:     u.DueDate,
	// 		// ApiTaskId:   u.ApiTaskId,
	// 		// Name:        u.Name,
	// 		// Persons:     u.Persons,
	// 		// RagRating:   u.RagRating,
	// 		// Status:      u.Status,
	// 		//Tasktype: u.Tasktype,
	// 	}

	// 	taskList = append(taskList, task)
	// }
	return allTaskList, err
}
