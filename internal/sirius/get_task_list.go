package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MetaData struct {
	TaskTypeCount []TypeAndCount `json:"taskTypeCount"`
}

type TypeAndCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type TaskList struct {
	Tasks      []model.Task          `json:"tasks"`
	Pages      model.PageInformation `json:"pages"`
	TotalTasks int                   `json:"total"`
	MetaData   MetaData              `json:"metadata"`
}

type TaskListParams struct {
	Team              model.Team
	Page              int
	PerPage           int
	TaskTypes         []model.TaskType
	SelectedTaskTypes []string
	Assignees         []string
	DueDateFrom       *time.Time
	DueDateTo         *time.Time
}

func (c *ApiClient) GetTaskList(ctx Context, params TaskListParams) (TaskList, error) {
	var v TaskList
	var teamIds []string

	if params.Team.Id != 0 {
		teamIds = []string{"teamIds[]=" + strconv.Itoa(params.Team.Id)}
	}
	for _, team := range params.Team.Teams {
		teamIds = append(teamIds, "teamIds[]="+strconv.Itoa(team.Id))
	}

	endpoint := fmt.Sprintf(
		"/api/v1/assignees/teams/tasks?%s&filter=%s&limit=%d&page=%d&sort=dueDate:asc",
		strings.Join(teamIds, "&"),
		params.CreateFilter(),
		params.PerPage,
		params.Page,
	)
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

func (p TaskListParams) CreateFilter() string {
	filter := "status:Not+started,"

	for _, t := range p.SelectedTaskTypes {
		if t == "ECM_TASKS" {
			p.SelectedTaskTypes = getEcmTaskTypesString(p.TaskTypes)
			break
		}
	}
	for _, t := range p.SelectedTaskTypes {
		filter += "type:" + t + ","
	}
	for _, a := range p.Assignees {
		filter += "assigneeid_or_null:" + a + ","
	}
	if p.DueDateFrom != nil {
		filter += "due_date_from:" + p.DueDateFrom.Format("2006-01-02") + ","
	}
	if p.DueDateTo != nil {
		filter += "due_date_to:" + p.DueDateTo.Format("2006-01-02") + ","
	}
	return strings.TrimRight(filter, ",")
}

func getEcmTaskTypesString(taskTypes []model.TaskType) []string {
	var ecmTasks []string
	for _, taskType := range taskTypes {
		if taskType.EcmTask {
			ecmTasks = append(ecmTasks, taskType.Handle)
		}
	}
	return ecmTasks
}
