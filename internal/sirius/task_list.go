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
	LimitedPagination []int
	FirstPage         int
	LastPage          int
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

func getPaginationLimits(TaskList TaskList, TaskDetails TaskDetails) []int {
	var twoBeforeCurrentPage int
	var twoAfterCurrentPage int
	if TaskList.Pages.PageCurrent > 2 {
		twoBeforeCurrentPage = TaskList.Pages.PageCurrent - 3
	} else {
		twoBeforeCurrentPage = 0
	}
	if TaskList.Pages.PageCurrent+2 <= TaskDetails.LastPage {
		twoAfterCurrentPage = TaskList.Pages.PageCurrent + 2
	} else if TaskList.Pages.PageCurrent+1 <= TaskDetails.LastPage {
		twoAfterCurrentPage = TaskList.Pages.PageCurrent + 1
	} else {
		twoAfterCurrentPage = TaskList.Pages.PageCurrent
	}
	return TaskDetails.ListOfPages[twoBeforeCurrentPage:twoAfterCurrentPage]
}

var teamID int
var taskType string

func (c *Client) GetTaskList(ctx Context, search int, displayTaskLimit int, selectedTeamMembers int, loggedInTeamId int, taskTypeSelected []string) (TaskList, TaskDetails, error) {
	var v TaskList
	var k TaskDetails
	var taskTypeFilters string

	if selectedTeamMembers == 0 {
		teamID = loggedInTeamId
	} else {
		teamID = selectedTeamMembers
	}

	if len(taskTypeSelected) == 0 {
		taskTypeFilters = ""
	} else if len(taskTypeSelected) == 1 {
		for _, s := range taskTypeSelected {
			taskTypeFilters += "type:" + s
		}
	} else if len(taskTypeSelected) > 1 {
		for _, s := range taskTypeSelected {
			taskTypeFilters += "type:" + s + ","
		}
		taskTypeFilterLength := len(taskTypeFilters)
		length := taskTypeFilterLength - 1
		taskTypeFilters = taskTypeFilters[0:length]
	}

	// req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/assignees/team/%d/tasks?limit=%d&page=%d&sort=dueDate:asc&filter=type:CWRD", teamID, displayTaskLimit, search), nil)
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/assignees/team/%d/tasks?filter=%s", teamID, taskTypeFilters), nil)

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

	if len(TaskDetails.ListOfPages) != 0 {
		TaskDetails.FirstPage = TaskDetails.ListOfPages[0]
		TaskDetails.LastPage = TaskDetails.ListOfPages[len(TaskDetails.ListOfPages)-1]
		TaskDetails.LimitedPagination = getPaginationLimits(TaskList, TaskDetails)
	} else {
		TaskDetails.FirstPage = 0
		TaskDetails.LastPage = 0
		TaskDetails.LimitedPagination = []int{0}
	}

	return TaskList, TaskDetails, err
}
