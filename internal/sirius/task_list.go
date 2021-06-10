package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AssigneeTeam struct {
	AssigneeTeamDisplayName string `json:"displayName"`
	AssigneeTeamId          int    `json:"id"`
}

type SupervisionTeam struct {
	SupervisionTeamDisplayName string `json:"displayName"`
	SupervisionTeamId          int    `json:"id"`
}

type SupervisionCaseOwner struct {
	SupervisionCaseOwnerName string            `json:"displayName"`
	SupervisionId            int               `json:"id"`
	SupervisionTeam          []SupervisionTeam `json:"teams"`
}

type CaseItemsDetails struct {
	CaseItemClient Clients `json:"client"`
}

type AssigneeDetails struct {
	AssigneeDisplayName string         `json:"displayName"`
	AssigneeId          int            `json:"id"`
	AssigneeTeams       []AssigneeTeam `json:"teams"`
}

type Clients struct {
	ClientId                   int                  `json:"id"`
	ClientCaseRecNumber        string               `json:"caseRecNumber"`
	ClientFirstName            string               `json:"firstname"`
	ClientSurname              string               `json:"surname"`
	ClientSupervisionCaseOwner SupervisionCaseOwner `json:"supervisionCaseOwner"`
}

type ApiTask struct {
	ApiTaskAssignee  AssigneeDetails    `json:"assignee"`
	ApiTaskCaseItems []CaseItemsDetails `json:"caseItems"`
	ApiClients       []Clients          `json:"clients"`
	ApiTaskDueDate   string             `json:"dueDate"`
	ApiTaskId        int                `json:"id"`
	ApiTaskHandle    string             `json:"type"`
	ApiTaskType      string             `json:"name"`
	ApiCaseOwnerTask bool               `json:"caseOwnerTask"`
	TaskTypeName     string
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
	LastFilter        string
	TaskTypeFilters   int
}

var teamID int

func (c *Client) GetTaskList(ctx Context, search int, displayTaskLimit int, selectedTeamMembers int, loggedInTeamId int, taskTypeSelected []string, LoadTasks []ApiTaskTypes) (TaskList, TaskDetails, error) {
	var v TaskList
	var k TaskDetails
	var taskTypeFilters string

	if selectedTeamMembers == 0 {
		teamID = loggedInTeamId
	} else {
		teamID = selectedTeamMembers
	}

	taskTypeFilters = createTaskTypeFilter(taskTypeSelected, taskTypeFilters)

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/assignees/team/%d/tasks?filter=status:Not+started,%s&limit=%d&page=%d&sort=dueDate:asc", teamID, taskTypeFilters, displayTaskLimit, search), nil)

	if err != nil {
		return v, k, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, k, err
	}
	// io.Copy(os.Stdout, resp.Body)

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

	TaskDetails.StoredTaskLimit = displayTaskLimit

	TaskDetails.ShowingUpperLimit = getShowingUpperLimitNumber(TaskList, displayTaskLimit)

	TaskDetails.ShowingLowerLimit = getShowingLowerLimitNumber(TaskList, displayTaskLimit)

	TaskDetails.LastFilter = getStoredTaskFilter(TaskDetails, taskTypeSelected, taskTypeFilters)

	TaskDetails.TaskTypeFilters = len(taskTypeSelected)

	if len(TaskDetails.ListOfPages) != 0 {
		TaskDetails.FirstPage = TaskDetails.ListOfPages[0]
		TaskDetails.LastPage = TaskDetails.ListOfPages[len(TaskDetails.ListOfPages)-1]
		TaskDetails.LimitedPagination = getPaginationLimits(TaskList, TaskDetails)
	} else {
		TaskDetails.FirstPage = 0
		TaskDetails.LastPage = 0
		TaskDetails.LimitedPagination = []int{0}
	}

	TaskList.WholeTaskList = setTaskTypeName(v.WholeTaskList, LoadTasks)
	fmt.Println(TaskList.WholeTaskList)

	return TaskList, TaskDetails, err
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

func getShowingLowerLimitNumber(TaskList TaskList, displayTaskLimit int) int {
	if TaskList.Pages.PageCurrent == 1 && TaskList.TotalTasks != 0 {
		return 1
	} else if TaskList.Pages.PageCurrent == 1 && TaskList.TotalTasks == 0 {
		return 0
	} else {
		previousPageNumber := TaskList.Pages.PageCurrent - 1
		return previousPageNumber*displayTaskLimit + 1
	}
}

func getShowingUpperLimitNumber(TaskList TaskList, displayTaskLimit int) int {
	if TaskList.Pages.PageCurrent*displayTaskLimit > TaskList.TotalTasks {
		return TaskList.TotalTasks
	} else {
		return TaskList.Pages.PageCurrent * displayTaskLimit
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

func createTaskTypeFilter(taskTypeSelected []string, taskTypeFilters string) string {
	if len(taskTypeSelected) == 1 {
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
	return taskTypeFilters
}

func getStoredTaskFilter(TaskDetails TaskDetails, taskTypeSelected []string, taskTypeFilters string) string {
	if TaskDetails.LastFilter == "" && len(taskTypeSelected) == 0 {
		return ""
	} else {
		return taskTypeFilters
	}
}

func setTaskTypeName(v []ApiTask, loadTasks []ApiTaskTypes) []ApiTask {
	var list []ApiTask
	for _, s := range v {
		task := ApiTask{
			ApiTaskAssignee: AssigneeDetails{
				AssigneeDisplayName: anotherFunc(s),
				AssigneeId:          anotherFuncId(s),
			},
			ApiTaskCaseItems: []CaseItemsDetails{
				CaseItemClient: []Clients{
					ClientId:                   anotherFuncClientId(s),
					ClientCaseRecNumber:        anotherFuncClient(s),
					ClientFirstName:            anotherFuncClient(s),
					ClientSurname:              anotherFuncClient(s),
					ClientSupervisionCaseOwner: anotherFuncClient(s),
				},
			},
			ApiTaskDueDate: s.ApiTaskDueDate,
			ApiTaskId:      s.ApiTaskId,
			ApiTaskHandle:  s.ApiTaskHandle,
			ApiTaskType:    s.ApiTaskType,
			TaskTypeName:   getTaskName(s, loadTasks),
		}
		list = append(list, task)
	}
	return list
}

func getTaskName(task ApiTask, loadTasks []ApiTaskTypes) string {
	for i := range loadTasks {
		if task.ApiTaskHandle == loadTasks[i].Handle {
			return loadTasks[i].Incomplete
		}
	}
	return task.ApiTaskType
}

func anotherFunc(s ApiTask) string {
	if s.ApiTaskAssignee.AssigneeDisplayName == "Unassigned" {
		if len(s.ApiClients) != 0 {
			return s.ApiClients[0].ClientSupervisionCaseOwner.SupervisionCaseOwnerName
		} else if len(s.ApiTaskCaseItems) != 0 {
			return s.ApiTaskCaseItems[0].CaseItemClient.ClientSupervisionCaseOwner.SupervisionCaseOwnerName
		}
	}
	return s.ApiTaskAssignee.AssigneeDisplayName
}

func anotherFuncId(s ApiTask) int {
	if s.ApiTaskAssignee.AssigneeId == 0 {
		if len(s.ApiClients) != 0 {
			return s.ApiClients[0].ClientSupervisionCaseOwner.SupervisionId
		} else if len(s.ApiTaskCaseItems) != 0 {
			return s.ApiTaskCaseItems[0].CaseItemClient.ClientSupervisionCaseOwner.SupervisionId
		}
	}
	return s.ApiTaskAssignee.AssigneeId
}
