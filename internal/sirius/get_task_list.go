package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CaseOwner struct {
	Name string     `json:"displayName"`
	Id   int        `json:"id"`
	Team []UserTeam `json:"teams"`
}

type UserTeam struct {
	Name string     `json:"displayName"`
	Id   int        `json:"id"`
	Team []UserTeam `json:"teams"`
}

type CaseItems struct {
	Client ApiClient `json:"client"`
}

type Deputy struct {
	Id          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Type        *struct {
		Handle string `json:"handle"`
	} `json:"deputyType"`
}

type ApiClient struct {
	Id            int       `json:"id"`
	CaseRecNumber string    `json:"caseRecNumber"`
	FirstName     string    `json:"firstname"`
	Surname       string    `json:"surname"`
	CaseOwner     CaseOwner `json:"supervisionCaseOwner"`
}

type ApiTask struct {
	Assignee          CaseOwner   `json:"assignee"`
	CaseItems         []CaseItems `json:"caseItems"`
	Clients           []ApiClient `json:"clients"`
	DueDate           string      `json:"dueDate"`
	Id                int         `json:"id"`
	TypeHandle        string      `json:"type"`
	TaskType          string      `json:"name"`
	TaskTypeName      string
	ClientInformation ApiClient
}

type PageInformation struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type MetaData struct {
	TaskTypeCount []TypeAndCount `json:"taskTypeCount"`
}

type TypeAndCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type TaskList struct {
	Tasks         []ApiTask       `json:"tasks"`
	Pages         PageInformation `json:"pages"`
	TotalTasks    int             `json:"total"`
	MetaData      MetaData        `json:"metadata"`
	ActiveFilters []string
}

func (c *Client) GetTaskList(ctx Context, search int, displayTaskLimit int, selectedTeam ReturnedTeamCollection, taskTypeSelected []string, LoadTasks []TaskTypes, selectedAssignees []string, dueDateFrom *time.Time, dueDateTo *time.Time) (TaskList, error) {
	var v TaskList
	var teamIds []string

	filter := CreateFilter(taskTypeSelected, selectedAssignees, dueDateFrom, dueDateTo)

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

	TaskList.Tasks = SetTaskTypeName(v.Tasks, LoadTasks)

	return TaskList, err
}

func (d *Deputy) GetURL() string {
	url := "/supervision/deputies/%d"
	if d.Type.Handle == "LAY" {
		url = "/supervision/#/deputy-hub/%d"
	}
	return fmt.Sprintf(url, d.Id)
}

func CreateFilter(taskTypeSelected []string, selectedAssignees []string, dueDateFrom *time.Time, dueDateTo *time.Time) string {
	filter := "status:Not+started,"
	for _, t := range taskTypeSelected {
		filter += "type:" + t + ","
	}
	for _, a := range selectedAssignees {
		filter += "assigneeid_or_null:" + a + ","
	}
	if dueDateFrom != nil {
		filter += "due_date_from:" + dueDateFrom.Format("2006-01-02") + ","
	}
	if dueDateTo != nil {
		filter += "due_date_to:" + dueDateTo.Format("2006-01-02") + ","
	}
	return strings.TrimRight(filter, ",")
}

func SetTaskTypeName(v []ApiTask, loadTasks []TaskTypes) []ApiTask {
	var list []ApiTask
	for _, s := range v {
		task := ApiTask{
			Assignee: CaseOwner{
				Name: GetAssigneeDisplayName(s),
				Id:   GetAssigneeId(s),
				Team: GetAssigneeTeams(s),
			},
			DueDate:           s.DueDate,
			Id:                s.Id,
			TypeHandle:        s.TypeHandle,
			TaskType:          s.TaskType,
			TaskTypeName:      GetTaskName(s, loadTasks),
			ClientInformation: GetClientInformation(s),
		}
		list = append(list, task)
	}
	return list
}

func GetTaskName(task ApiTask, loadTasks []TaskTypes) string {
	for i := range loadTasks {
		if task.TypeHandle == loadTasks[i].Handle {
			return loadTasks[i].Name
		}
	}
	return task.TaskType
}

func GetAssigneeDisplayName(s ApiTask) string {
	if s.Assignee.Name == "Unassigned" {
		if len(s.Clients) != 0 {
			return s.Clients[0].CaseOwner.Name
		} else if len(s.CaseItems) != 0 {
			return s.CaseItems[0].Client.CaseOwner.Name
		}
	}
	return s.Assignee.Name
}

func GetAssigneeTeams(s ApiTask) []UserTeam {
	if len(s.Assignee.Team) == 0 {
		if len(s.Clients) != 0 {
			return s.Clients[0].CaseOwner.Team
		} else if len(s.CaseItems) != 0 {
			return s.CaseItems[0].Client.CaseOwner.Team
		}
	}
	return s.Assignee.Team
}

func GetAssigneeId(s ApiTask) int {
	if s.Assignee.Id == 0 {
		if len(s.Clients) != 0 {
			return s.Clients[0].CaseOwner.Id
		} else if len(s.CaseItems) != 0 {
			return s.CaseItems[0].Client.CaseOwner.Id
		}
	}
	return s.Assignee.Id
}

func GetClientInformation(s ApiTask) ApiClient {
	if len(s.CaseItems) != 0 {
		return s.CaseItems[0].Client
	}
	return s.Clients[0]
}
