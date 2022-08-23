package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CaseManagement struct {
	CaseManagerName string     `json:"displayName"`
	Id              int        `json:"id"`
	Team            []UserTeam `json:"teams"`
}

type UserTeam struct {
	Name string     `json:"displayName"`
	Id   int        `json:"id"`
	Team []UserTeam `json:"teams"`
}

type CaseItemsDetails struct {
	CaseItemClient Clients `json:"client"`
}

type Clients struct {
	ClientId                   int            `json:"id"`
	ClientCaseRecNumber        string         `json:"caseRecNumber"`
	ClientFirstName            string         `json:"firstname"`
	ClientSurname              string         `json:"surname"`
	ClientSupervisionCaseOwner CaseManagement `json:"supervisionCaseOwner"`
}

type ApiTask struct {
	ApiTaskAssignee   CaseManagement     `json:"assignee"`
	ApiTaskCaseItems  []CaseItemsDetails `json:"caseItems"`
	ApiClients        []Clients          `json:"clients"`
	ApiTaskDueDate    string             `json:"dueDate"`
	ApiTaskId         int                `json:"id"`
	ApiTaskHandle     string             `json:"type"`
	ApiTaskType       string             `json:"name"`
	ApiCaseOwnerTask  bool               `json:"caseOwnerTask"`
	TaskTypeName      string
	ClientInformation Clients
}

type PageDetails struct {
	PageCurrent int `json:"current"`
	PageTotal   int `json:"total"`
}

type TaskList struct {
	WholeTaskList []ApiTask   `json:"tasks"`
	Pages         PageDetails `json:"pages"`
	TotalTasks    int         `json:"total"`
	ActiveFilters []string
}

var teamID int

func (c *Client) GetTaskList(ctx Context, search int, displayTaskLimit int, selectedTeamId int, loggedInTeamId int, taskTypeSelected []string, LoadTasks []ApiTaskTypes, assigneeSelected []string) (TaskList, int, error) {
	var v TaskList
	var taskTypeFilters string
	var assigneeFilters string

	if selectedTeamId == 0 {
		teamID = loggedInTeamId
	} else {
		teamID = selectedTeamId
	}

	taskTypeFilters = CreateTaskTypeFilter(taskTypeSelected, taskTypeFilters)
	assigneeFilters = CreateAssigneeFilter(assigneeSelected, assigneeFilters)
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/assignees/team/%d/tasks?filter=status:Not+started,%s%s&limit=%d&page=%d&sort=dueDate:asc", teamID, taskTypeFilters, assigneeFilters, displayTaskLimit, search), nil)

	if err != nil {
		return v, 0, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, 0, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, 0, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, 0, err
	}

	TaskList := v

	TaskList.WholeTaskList = SetTaskTypeName(v.WholeTaskList, LoadTasks)

	return TaskList, teamID, err
}

func CreateTaskTypeFilter(taskTypeSelected []string, taskTypeFilters string) string {
	if len(taskTypeSelected) == 0 {
		taskTypeFilters += ","
	} else {
		for _, s := range taskTypeSelected {
			taskTypeFilters += "type:" + s + ","
		}
	}
	return taskTypeFilters
}

func CreateAssigneeFilter(assigneeSelected []string, assigneeFilters string) string {
	if len(assigneeSelected) == 1 {
		for _, s := range assigneeSelected {
			assigneeFilters += "assigneeid_or_null:" + s
		}
	} else if len(assigneeSelected) > 1 {
		for _, s := range assigneeSelected {
			assigneeFilters += "assigneeid_or_null:" + s + ","
		}
		assigneeFilterLength := len(assigneeFilters)
		length := assigneeFilterLength - 1
		assigneeFilters = assigneeFilters[0:length]
	}
	return assigneeFilters
}

func SetTaskTypeName(v []ApiTask, loadTasks []ApiTaskTypes) []ApiTask {
	var list []ApiTask
	for _, s := range v {
		task := ApiTask{
			ApiTaskAssignee: CaseManagement{
				CaseManagerName: GetAssigneeDisplayName(s),
				Id:              GetAssigneeId(s),
				Team:            GetAssigneeTeams(s),
			},
			ApiTaskDueDate:    s.ApiTaskDueDate,
			ApiTaskId:         s.ApiTaskId,
			ApiTaskHandle:     s.ApiTaskHandle,
			ApiTaskType:       s.ApiTaskType,
			TaskTypeName:      GetTaskName(s, loadTasks),
			ClientInformation: GetClientInformation(s),
		}
		list = append(list, task)
	}
	return list
}

func GetTaskName(task ApiTask, loadTasks []ApiTaskTypes) string {
	for i := range loadTasks {
		if task.ApiTaskHandle == loadTasks[i].Handle {
			return loadTasks[i].Incomplete
		}
	}
	return task.ApiTaskType
}

func GetAssigneeDisplayName(s ApiTask) string {
	if s.ApiTaskAssignee.CaseManagerName == "Unassigned" {
		if len(s.ApiClients) != 0 {
			return s.ApiClients[0].ClientSupervisionCaseOwner.CaseManagerName
		} else if len(s.ApiTaskCaseItems) != 0 {
			return s.ApiTaskCaseItems[0].CaseItemClient.ClientSupervisionCaseOwner.CaseManagerName
		}
	}
	return s.ApiTaskAssignee.CaseManagerName
}

func GetAssigneeTeams(s ApiTask) []UserTeam {
	if len(s.ApiTaskAssignee.Team) == 0 {
		if len(s.ApiClients) != 0 {
			return s.ApiClients[0].ClientSupervisionCaseOwner.Team
		} else if len(s.ApiTaskCaseItems) != 0 {
			return s.ApiTaskCaseItems[0].CaseItemClient.ClientSupervisionCaseOwner.Team
		}
	}
	return s.ApiTaskAssignee.Team
}

func GetAssigneeId(s ApiTask) int {
	if s.ApiTaskAssignee.Id == 0 {
		if len(s.ApiClients) != 0 {
			return s.ApiClients[0].ClientSupervisionCaseOwner.Id
		} else if len(s.ApiTaskCaseItems) != 0 {
			return s.ApiTaskCaseItems[0].CaseItemClient.ClientSupervisionCaseOwner.Id
		}
	}
	return s.ApiTaskAssignee.Id
}

func GetClientInformation(s ApiTask) Clients {
	if len(s.ApiTaskCaseItems) != 0 {
		return s.ApiTaskCaseItems[0].CaseItemClient
	}
	return s.ApiClients[0]
}
