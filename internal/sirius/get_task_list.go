package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RefData struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type Assignee struct {
	Name  string     `json:"displayName"`
	Id    int        `json:"id"`
	Teams []UserTeam `json:"teams"`
}

type UserTeam struct {
	Name string     `json:"displayName"`
	Id   int        `json:"id"`
	Team []UserTeam `json:"teams"`
}

type CaseItem struct {
	Client Client `json:"client"`
}

type Deputy struct {
	Id          int     `json:"id"`
	DisplayName string  `json:"displayName"`
	Type        RefData `json:"deputyType"`
}

type Client struct {
	Id                   int      `json:"id"`
	CaseRecNumber        string   `json:"caseRecNumber"`
	FirstName            string   `json:"firstname"`
	Surname              string   `json:"surname"`
	SupervisionCaseOwner Assignee `json:"supervisionCaseOwner"`
	FeePayer             Deputy   `json:"feePayer"`
	Orders               []Order  `json:"cases"`
	SupervisionLevel     string   `json:"supervisionLevel"`
}

type Task struct {
	Assignee      Assignee   `json:"assignee"`
	CaseItems     []CaseItem `json:"caseItems"`
	Clients       []Client   `json:"clients"`
	DueDate       string     `json:"dueDate"`
	Id            int        `json:"id"`
	Type          string     `json:"type"`
	Name          string     `json:"name"`
	CaseOwnerTask bool       `json:"caseOwnerTask"`
	IsPriority    bool       `json:"isPriority"`
}

type DueDateStatus struct {
	Name   string
	Colour string
}

type PageInformation struct {
	PageCurrent int `json:"current"`
	PageTotal   int `json:"total"`
}

type MetaData struct {
	TaskTypeCount []TypeAndCount `json:"taskTypeCount"`
}

type TypeAndCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type TaskList struct {
	Tasks      []Task          `json:"tasks"`
	Pages      PageInformation `json:"pages"`
	TotalTasks int             `json:"total"`
	MetaData   MetaData        `json:"metadata"`
}

func (c *ApiClient) GetTaskList(ctx Context, search int, displayTaskLimit int, selectedTeam Team, taskTypeSelected []string, taskTypes []TaskType, selectedAssignees []string, dueDateFrom *time.Time, dueDateTo *time.Time) (TaskList, error) {
	var v TaskList
	var teamIds []string

	filter := CreateFilter(taskTypeSelected, selectedAssignees, dueDateFrom, dueDateTo, taskTypes)

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

	return v, nil
}

func (d Deputy) GetURL() string {
	url := "/supervision/deputies/%d"
	if d.Type.Handle == "LAY" {
		url = "/supervision/#/deputy-hub/%d"
	}
	return fmt.Sprintf(url, d.Id)
}

func CreateFilter(taskTypeSelected []string, selectedAssignees []string, dueDateFrom *time.Time, dueDateTo *time.Time, taskTypes []TaskType) string {
	filter := "status:Not+started,"

	for _, t := range taskTypeSelected {
		if t == "ECM_TASKS" {
			taskTypeSelected = getEcmTaskTypesString(taskTypes)
			break
		}
	}
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

func (t Task) GetDueDateStatus(now ...time.Time) DueDateStatus {
	removeTime := func(t time.Time) time.Time {
		y, m, d := t.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	}

	today := removeTime(time.Now())
	if len(now) > 0 {
		today = removeTime(now[0])
	}

	dueDate, _ := time.Parse("02/01/2006", t.DueDate)
	dueDate = removeTime(dueDate)

	daysUntilNextWeek := int((7 + (time.Monday - today.Weekday())) % 7)
	startOfNextWeek := today.AddDate(0, 0, daysUntilNextWeek)
	endOfNextWeek := startOfNextWeek.AddDate(0, 0, 6)

	if dueDate.Before(today) {
		return DueDateStatus{"Overdue", "red"}
	} else if dueDate.Equal(today) {
		return DueDateStatus{"Due Today", "red"}
	} else if dueDate.Sub(today).Hours() == 24 && dueDate.Before(startOfNextWeek) {
		return DueDateStatus{"Due Tomorrow", "orange"}
	} else if dueDate.Before(startOfNextWeek) {
		return DueDateStatus{"Due This Week", "orange"}
	} else if dueDate.Before(endOfNextWeek) || dueDate.Equal(endOfNextWeek) {
		return DueDateStatus{"Due Next Week", "green"}
	}

	return DueDateStatus{"", ""}
}

func (t Task) GetClient() Client {
	if len(t.CaseItems) != 0 {
		return t.CaseItems[0].Client
	}
	return t.Clients[0]
}

func (t Task) GetName(taskTypes []TaskType) string {
	for _, taskType := range taskTypes {
		if t.Type == taskType.Handle {
			return taskType.Incomplete
		}
	}
	return t.Name
}

func (t Task) GetAssignee() Assignee {
	if t.Assignee.Name == "Unassigned" {
		if len(t.Clients) != 0 {
			return t.Clients[0].SupervisionCaseOwner
		} else if len(t.CaseItems) != 0 {
			return t.CaseItems[0].Client.SupervisionCaseOwner
		}
	}
	return t.Assignee
}

func (c Client) GetReportDueDate() string {
	if len(c.Orders) > 0 {
		return c.Orders[0].LatestAnnualReport.DueDate
	}
	return ""
}

func (c Client) GetStatus() string {
	orderStatuses := make(map[string]string)

	for _, order := range c.Orders {
		label := order.Status.Label
		orderStatuses[label] = label
	}

	statuses := []string{"Active", "Open", "Closed", "Duplicate"}
	for _, status := range statuses {
		if _, found := orderStatuses[status]; found {
			return status
		}
	}

	return ""
}

func getEcmTaskTypesString(taskTypes []TaskType) []string {
	var ecmTasks []string
	for _, taskType := range taskTypes {
		if taskType.EcmTask {
			ecmTasks = append(ecmTasks, taskType.Handle)
		}
	}
	return ecmTasks
}
