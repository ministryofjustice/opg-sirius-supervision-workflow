package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SupervisionCaseOwnerDetail struct {
	SupervisionCaseOwnerName string `json:"displayName"`
}

type ClientDetails struct {
	ClientCaseRecNumber        string                     `json:"caseRecNumber"`
	ClientFirstName            string                     `json:"firstname"`
	ClientId                   int                        `json:"id"`
	ClientSupervisionCaseOwner SupervisionCaseOwnerDetail `json:"supervisionCaseOwner"`
	ClientSurname              string                     `json:"surname"`
}

type CaseItemsDetails struct {
	CaseItemClient ClientDetails `json:"client"`
}

type AssigneeDetails struct {
	AssigneeDisplayName string `json:"displayName"`
	AssigneeId          int    `json:"id"`
}

type ApiTask struct {
	ApiTaskAssignee  AssigneeDetails    `json:"assignee"`
	ApiTaskCaseItems []CaseItemsDetails `json:"caseItems"`
	ApiTaskDueDate   string             `json:"dueDate"`
	ApiTaskId        int                `json:"id"`
	ApiTaskType      string             `json:"name"`
}

type PageDetails struct {
	PageCurrent int `json:"current"`
	PageTotal   int `json:"total"`
}

type TaskList struct {
	WholeTaskList []ApiTask   `json:"tasks"`
	Pages         PageDetails `json:"pages"`
	TotalTasks    int         `json:"total"`
}

type TaskDetails struct {
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

func getNextPageNumber(TaskList TaskList, search int) int {
	if search < TaskList.Pages.PageTotal {
		if search == 0 {
			return search + 2
		} else {
			return search + 1
		}
	} else {
		return TaskList.Pages.PageTotal
	}
}

func getStoredTaskLimitNumber(TaskDetails TaskDetails, displayTaskLimit int) int {
	if TaskDetails.StoredTaskLimit == 0 && displayTaskLimit == 0 {
		return 25
	} else {
		return displayTaskLimit
	}
}

func getShowingLowerLimitNumber(TaskList TaskList, TaskDetails TaskDetails) int {
	if TaskList.Pages.PageCurrent == 1 && TaskList.TotalTasks != 0 {
		return 1
	} else if TaskList.Pages.PageCurrent == 1 && TaskList.TotalTasks == 0 {
		return 0
	} else {
		previousPageNumber := TaskList.Pages.PageCurrent - 1
		return previousPageNumber*TaskDetails.StoredTaskLimit + 1
	}
}

func getShowingUpperLimitNumber(TaskList TaskList, TaskDetails TaskDetails) int {
	if TaskList.Pages.PageCurrent*TaskDetails.StoredTaskLimit > TaskList.TotalTasks {
		return TaskList.TotalTasks
	} else {
		return TaskList.Pages.PageCurrent * TaskDetails.StoredTaskLimit
	}
}

var teamID int

func (c *Client) GetTaskList(ctx Context, search int, displayTaskLimit int, selectedTeamMembers int, loggedInTeamId int) (TaskList, TaskDetails, error) {
	var v TaskList
	var k TaskDetails

	if selectedTeamMembers == 0 {
		teamID = loggedInTeamId
	} else {
		teamID = selectedTeamMembers
	}

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/assignees/team/%d/tasks?limit=%d&page=%d&sort=dueDate:asc", teamID, displayTaskLimit, search), nil)
	// req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/assignees/team/%d/tasks", teamID), nil)
	if err != nil {
		return v, k, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, k, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, k, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, k, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, k, err
	}

	TaskList := v
	TaskDetails := k

	for i := 1; i < TaskList.Pages.PageTotal+1; i++ {
		TaskDetails.ListOfPages = append(TaskDetails.ListOfPages, i)
	}

	TaskDetails.PreviousPage = getPreviousPageNumber(search)

	TaskDetails.NextPage = getNextPageNumber(TaskList, search)

	TaskDetails.StoredTaskLimit = getStoredTaskLimitNumber(TaskDetails, displayTaskLimit)

	TaskDetails.ShowingUpperLimit = getShowingUpperLimitNumber(TaskList, TaskDetails)

	TaskDetails.ShowingLowerLimit = getShowingLowerLimitNumber(TaskList, TaskDetails)

	return TaskList, TaskDetails, err
}
