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
	TaskTypeCount []TypeAndCount     `json:"taskTypeCount"`
	AssigneeCount []AssigneeAndCount `json:"assigneeTaskCount"`
}

type TypeAndCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type AssigneeAndCount struct {
	AssigneeId int `json:"assignee"`
	Count      int `json:"count"`
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
	TaskTypeCategory  string
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
		"/api/v1/assignees/teams/tasks?%s&filter=%s&limit=%d&page=%d&sort=%s",
		strings.Join(teamIds, "&"),
		params.CreateFilter(),
		params.PerPage,
		params.Page,
		"ispriority:desc,duedate:asc,id:asc",
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

	fmt.Println("v")
	fmt.Println(v.MetaData)
	fmt.Println(v.MetaData.AssigneeCount)

	return v, nil
}

func (p TaskListParams) CreateFilter() string {
	filter := "status:Not+started,"

	for _, t := range p.SelectedTaskTypes {
		if t == TaskTypeEcmHandle {
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
	if p.TaskTypeCategory != "" {
		filter += "task_type_category:" + p.TaskTypeCategory
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

func (tl TaskList) CalculateTaskTypeCounts(taskTypes []model.TaskType) []model.TaskType {
	ecmTasksCount := 0
	getTaskTypeCount := func(taskType string) int {
		for _, q := range tl.MetaData.TaskTypeCount {
			if taskType == q.Type {
				return q.Count
			}
		}
		return 0
	}

	for i, t := range taskTypes {
		taskTypes[i].TaskCount = getTaskTypeCount(t.Handle)
		if t.EcmTask {
			ecmTasksCount += taskTypes[i].TaskCount
		}
	}

	if ecmTasksCount > 0 {
		for i, t := range taskTypes {
			if t.Handle == TaskTypeEcmHandle {
				taskTypes[i].TaskCount = ecmTasksCount
			}
		}
	}

	return taskTypes
}

//func (tl TaskList) CalculateAssigneeTaskCounts(assignees []model.Assignee) {
//	fmt.Println("assignees")
//	fmt.Println(assignees)
//
//	fmt.Println("meta")
//	fmt.Println(tl.MetaData.AssigneeCount)
//	for i, t := range assignees {
//		for assigneeId, count := range tl.MetaData.AssigneeCount {
//			if i == assigneeId {
//				assignees.
//			}
//		}
//	}
//	//
//	//for i, t := range taskTypes {
//	//	taskTypes[i].TaskCount = getTaskTypeCount(t.Handle)
//	//	if t.EcmTask {
//	//		ecmTasksCount += taskTypes[i].TaskCount
//	//	}
//	//}
//	//
//	//if ecmTasksCount > 0 {
//	//	for i, t := range taskTypes {
//	//		if t.Handle == TaskTypeEcmHandle {
//	//			taskTypes[i].TaskCount = ecmTasksCount
//	//		}
//	//	}
//	//}
//
//	return taskTypes
//}
