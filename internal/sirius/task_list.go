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
}

var teamID int

func (c *Client) GetTaskList(ctx Context, search int, displayTaskLimit int, selectedTeamMembers int, loggedInTeamId int, taskTypeSelected []string, LoadTasks []ApiTaskTypes) (TaskList, error) {
	var v TaskList
	var taskTypeFilters string

	if selectedTeamMembers == 0 {
		teamID = loggedInTeamId
	} else {
		teamID = selectedTeamMembers
	}

	taskTypeFilters = createTaskTypeFilter(taskTypeSelected, taskTypeFilters)

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/assignees/team/%d/tasks?filter=status:Not+started,%s&limit=%d&page=%d&sort=dueDate:asc", teamID, taskTypeFilters, displayTaskLimit, search), nil)

	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	// io.Copy(os.Stdout, resp.Body)

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

	TaskList.WholeTaskList = setTaskTypeName(v.WholeTaskList, LoadTasks)

	return TaskList, err
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
			ApiTaskAssignee: CaseManagement{
				CaseManagerName: getAssigneeDisplayName(s),
				Id:              getAssigneeId(s),
				Team:            getAssigneeTeams(s),
			},
			ApiTaskDueDate:    s.ApiTaskDueDate,
			ApiTaskId:         s.ApiTaskId,
			ApiTaskHandle:     s.ApiTaskHandle,
			ApiTaskType:       s.ApiTaskType,
			TaskTypeName:      getTaskName(s, loadTasks),
			ClientInformation: getClientInformation(s),
		}
		fmt.Println("new task")
		fmt.Println(task)
		list = append(list, task)

	}
	fmt.Println("list")
	fmt.Println(list)
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

func getAssigneeDisplayName(s ApiTask) string {
	if s.ApiTaskAssignee.CaseManagerName == "Unassigned" {
		if len(s.ApiClients) != 0 {
			return s.ApiClients[0].ClientSupervisionCaseOwner.CaseManagerName
		} else if len(s.ApiTaskCaseItems) != 0 {
			return s.ApiTaskCaseItems[0].CaseItemClient.ClientSupervisionCaseOwner.CaseManagerName
		}
	}
	return s.ApiTaskAssignee.CaseManagerName
}

func getAssigneeTeams(s ApiTask) []UserTeam {
	if len(s.ApiTaskAssignee.Team) == 0 {
		if len(s.ApiClients) != 0 {
			return s.ApiClients[0].ClientSupervisionCaseOwner.Team
		} else if len(s.ApiTaskCaseItems) != 0 {
			return s.ApiTaskCaseItems[0].CaseItemClient.ClientSupervisionCaseOwner.Team
		}
	}
	return s.ApiTaskAssignee.Team
}

func getAssigneeId(s ApiTask) int {
	if s.ApiTaskAssignee.Id == 0 {
		if len(s.ApiClients) != 0 {
			return s.ApiClients[0].ClientSupervisionCaseOwner.Id
		} else if len(s.ApiTaskCaseItems) != 0 {
			return s.ApiTaskCaseItems[0].CaseItemClient.ClientSupervisionCaseOwner.Id
		}
	}
	return s.ApiTaskAssignee.Id
}

func getClientInformation(s ApiTask) Clients {
	if len(s.ApiTaskCaseItems) != 0 {
		return s.ApiTaskCaseItems[0].CaseItemClient
	}
	return s.ApiClients[0]
}
