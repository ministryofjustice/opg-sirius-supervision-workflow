package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type RefData struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

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

type Deputy struct {
	Id          int     `json:"id"`
	DisplayName string  `json:"displayName"`
	Type        RefData `json:"deputyType"`
}

type Clients struct {
	ClientId                   int            `json:"id"`
	ClientCaseRecNumber        string         `json:"caseRecNumber"`
	ClientFirstName            string         `json:"firstname"`
	ClientSurname              string         `json:"surname"`
	ClientSupervisionCaseOwner CaseManagement `json:"supervisionCaseOwner"`
	FeePayer                   Deputy         `json:"feePayer"`
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

type PageInformation struct {
	PageCurrent int `json:"current"`
	PageTotal   int `json:"total"`
}

type TaskList struct {
	WholeTaskList []ApiTask       `json:"tasks"`
	Pages         PageInformation `json:"pages"`
	TotalTasks    int             `json:"total"`
	ActiveFilters []string
}

func (c *Client) GetTaskList(ctx Context, search int, displayTaskLimit int, selectedTeam ReturnedTeamCollection, taskTypeSelected []string, LoadTasks []ApiTaskTypes, selectedAssignees []string) (TaskList, error) {
	var v TaskList
	var teamIds []string

	filter := CreateFilter(taskTypeSelected, selectedAssignees)

	if selectedTeam.Id != 0 {
		teamIds = []string{"teamIds[]=" + strconv.Itoa(selectedTeam.Id)}
	}
	for _, team := range selectedTeam.Teams {
		teamIds = append(teamIds, "teamIds[]="+strconv.Itoa(team.Id))
	}

	endpoint := fmt.Sprintf("/api/v1/assignees/teams/tasks?%s&filter=%s&limit=%d&page=%d&sort=dueDate:asc", strings.Join(teamIds, "&"), filter, displayTaskLimit, search)
	req, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)

	if err != nil {
		c.logErrorRequest(req, err)
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.logResponse(req, resp, err)
		return v, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.logResponse(req, resp, err)
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logResponse(req, resp, err)
		return v, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		c.logResponse(req, resp, err)
		return v, err
	}

	TaskList := v

	TaskList.WholeTaskList = SetTaskTypeName(v.WholeTaskList, LoadTasks)

	return TaskList, err
}

func (d *Deputy) GetURL() string {
	url := "/supervision/deputies/%d"
	if d.Type.Handle == "LAY" {
		url = "/supervision/#/deputy-hub/%d"
	}
	return fmt.Sprintf(url, d.Id)
}

func CreateFilter(taskTypeSelected []string, selectedAssignees []string) string {
	filter := "status:Not+started,"
	for _, t := range taskTypeSelected {
		filter += "type:" + t + ","
	}
	for _, a := range selectedAssignees {
		filter += "assigneeid_or_null:" + a + ","
	}
	return strings.TrimRight(filter, ",")
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
