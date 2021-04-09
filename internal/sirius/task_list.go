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

type PageDetails struct {
	PageCurrent int `json:"current"`
	PageTotal   int `json:"total"`
}

type TaskList struct {
	WholeTaskList     []ApiTask   `json:"tasks"`
	Pages             PageDetails `json:"pages"`
	TotalTasks        int         `json:"total"`
	ListOfPages       []int
	PreviousPage      int
	NextPage          int
	StoredTaskLimit   int
	ShowingUpperLimit int
	ShowingLowerLimit int
}

func getPreviousPageNumber(search int) int {
	if search <= 1 {
		return 1
	} else {
		return search - 1
	}
}

func (c *Client) GetTaskList(ctx Context, search int, displayTaskLimit int) (TaskList, error) {
	var v TaskList

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/assignees/team/tasks?limit=%d&page=%d&sort=dueDate:asc", displayTaskLimit, search), nil)
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

	for i := 1; i < TaskList.Pages.PageTotal+1; i++ {
		TaskList.ListOfPages = append(TaskList.ListOfPages, i)
	}

	TaskList.PreviousPage = getPreviousPageNumber(search)

	if search < TaskList.Pages.PageTotal {
		if search == 0 {
			TaskList.NextPage = (search + 2)
		} else {
			TaskList.NextPage = (search + 1)
		}
	} else {
		TaskList.NextPage = TaskList.Pages.PageTotal
	}
	if TaskList.StoredTaskLimit == 0 && displayTaskLimit == 0 {
		TaskList.StoredTaskLimit = 25
	} else {
		TaskList.StoredTaskLimit = displayTaskLimit
	}

	TaskList.ShowingUpperLimit = TaskList.Pages.PageCurrent * TaskList.StoredTaskLimit

	if TaskList.Pages.PageCurrent == 1 {
		TaskList.ShowingLowerLimit = 1
	} else {
		TaskList.ShowingLowerLimit = ((TaskList.Pages.PageCurrent - 1) * TaskList.StoredTaskLimit) + 1
	}

	if TaskList.ShowingUpperLimit > TaskList.TotalTasks {
		TaskList.ShowingUpperLimit = TaskList.TotalTasks
	}

	return TaskList, err
}
