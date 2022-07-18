package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CaseManagement struct {
	Name  string     `json:"displayName"`
	Id    int        `json:"id"`
	Teams []UserTeam `json:"teams"`
}

type UserTeam struct {
	Name  string     `json:"displayName"`
	Id    int        `json:"id"`
	Teams []UserTeam `json:"teams"`
}

type CaseItemsDetails struct {
	Client SupervisionClient `json:"client"`
}

type SupervisionClient struct {
	Id                   int            `json:"id"`
	CaseRecNumber        string         `json:"caseRecNumber"`
	FirstName            string         `json:"firstname"`
	Surname              string         `json:"surname"`
	SupervisionCaseOwner CaseManagement `json:"supervisionCaseOwner"`
}

type Task struct {
	Id                int                 `json:"id"`
	Assignee          CaseManagement      `json:"assignee"`
	CaseItems         []CaseItemsDetails  `json:"caseItems"`
	Clients           []SupervisionClient `json:"clients"`
	DueDate           string              `json:"dueDate"`
	Handle            string              `json:"type"`
	Type              string              `json:"name"`
	CaseOwnerTask     bool                `json:"caseOwnerTask"`
	Name              string
	ClientInformation SupervisionClient
}

type PageDetails struct {
	PageCurrent int `json:"current"`
	PageTotal   int `json:"total"`
}

type TaskList struct {
	WholeTaskList []Task      `json:"tasks"`
	Pages         PageDetails `json:"pages"`
	TotalTasks    int         `json:"total"`
	ActiveFilters []string
}

var teamID int

func (c *Client) GetTaskList(ctx Context, search int, displayTaskLimit int, selectedTeamId int, loggedInTeamId int, taskTypeSelected []string, LoadTasks []TaskType, assigneeSelected []string) (TaskList, int, error) {
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
	} else if len(taskTypeSelected) == 1 {
		for _, s := range taskTypeSelected {
			taskTypeFilters += "type:" + s + ","
		}
	} else if len(taskTypeSelected) > 1 {
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

func SetTaskTypeName(v []Task, loadTasks []TaskType) []Task {
	var list []Task
	for _, s := range v {
		task := Task{
			Assignee: CaseManagement{
				Name:  GetAssigneeDisplayName(s),
				Id:    GetAssigneeId(s),
				Teams: GetAssigneeTeams(s),
			},
			DueDate:           s.DueDate,
			Id:                s.Id,
			Handle:            s.Handle,
			Type:              s.Type,
			Name:              GetTaskName(s, loadTasks),
			ClientInformation: GetClientInformation(s),
		}
		list = append(list, task)
	}
	return list
}

func GetTaskName(task Task, loadTasks []TaskType) string {
	for i := range loadTasks {
		if task.Handle == loadTasks[i].Handle {
			return loadTasks[i].Incomplete
		}
	}
	return task.Type
}

func GetAssigneeDisplayName(s Task) string {
	if s.Assignee.Name == "Unassigned" {
		if len(s.Clients) != 0 {
			return s.Clients[0].SupervisionCaseOwner.Name
		} else if len(s.CaseItems) != 0 {
			return s.CaseItems[0].Client.SupervisionCaseOwner.Name
		}
	}
	return s.Assignee.Name
}

func GetAssigneeTeams(s Task) []UserTeam {
	if len(s.Assignee.Teams) == 0 {
		if len(s.Clients) != 0 {
			return s.Clients[0].SupervisionCaseOwner.Teams
		} else if len(s.CaseItems) != 0 {
			return s.CaseItems[0].Client.SupervisionCaseOwner.Teams
		}
	}
	return s.Assignee.Teams
}

func GetAssigneeId(s Task) int {
	if s.Assignee.Id == 0 {
		if len(s.Clients) != 0 {
			return s.Clients[0].SupervisionCaseOwner.Id
		} else if len(s.CaseItems) != 0 {
			return s.CaseItems[0].Client.SupervisionCaseOwner.Id
		}
	}
	return s.Assignee.Id
}

func GetClientInformation(s Task) SupervisionClient {
	if len(s.CaseItems) != 0 {
		return s.CaseItems[0].Client
	}
	return s.Clients[0]
}
