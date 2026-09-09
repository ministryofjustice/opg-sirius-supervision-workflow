package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
)

type TaskMetaData struct {
	TaskTypeCount []TypeAndCount           `json:"taskTypeCount"`
	AssigneeCount []model.AssigneeAndCount `json:"assigneeTaskCount"`
}

type TypeAndCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type TaskList struct {
	Tasks      []model.Task          `json:"tasks"`
	Pages      model.PageInformation `json:"pages"`
	TotalTasks int                   `json:"total"`
	MetaData   TaskMetaData          `json:"metadata"`
}

type TaskListParams struct {
	Team              model.Team
	Page              int
	PerPage           int
	TaskTypes         []model.TaskType
	TaskTypeCategory  string
	SelectedTaskTypes []string
	Assignees         []string
	Deputies          []string
	DueDateFrom       *time.Time
	DueDateTo         *time.Time
}

type taskMetaDataResponse struct {
	TaskTypeCount []TypeAndCount           `json:"taskTypeCount"`
	AssigneeCount []model.AssigneeAndCount `json:"assigneeTaskCount"`
}

type taskResponse struct {
	Assignee      *assigneeWithTeamsResponse `json:"assignee,omitempty"`
	Orders        []orderResponse            `json:"caseItems,omitempty"`
	Persons       []clientSummaryResponse    `json:"persons,omitempty"`
	Clients       []clientSummaryResponse    `json:"clients,omitempty"`
	Deputies      []deputyResponse           `json:"deputies,omitempty"`
	Status        string                     `json:"status,omitempty"`
	DueDate       string                     `json:"dueDate"`
	Id            int                        `json:"id"`
	Type          string                     `json:"type"`
	Name          string                     `json:"name"`
	Description   string                     `json:"description,omitempty"`
	RAGRating     int                        `json:"ragRating,omitempty"`
	CreatedTime   string                     `json:"createdTime,omitempty"`
	CaseOwnerTask bool                       `json:"caseOwnerTask"`
	IsPriority    bool                       `json:"isPriority"`
}

func (r taskResponse) toTask() (model.Task, error) {
	var orders []model.Order
	if r.Orders != nil {
		orders = make([]model.Order, 0, len(r.Orders))
		for _, order := range r.Orders {
			mappedOrder, err := order.toOrder()
			if err != nil {
				return model.Task{}, err
			}
			orders = append(orders, mappedOrder)
		}
	}

	var clients []model.Client
	if r.Clients != nil {
		clients = make([]model.Client, 0, len(r.Clients))
		for _, client := range r.Clients {
			client := client
			clients = append(clients, client.toClient())
		}
	}

	var deputies []model.Deputy
	if r.Deputies != nil {
		deputies = make([]model.Deputy, 0, len(r.Deputies))
		for _, deputy := range r.Deputies {
			mappedDeputy, err := deputy.toDeputy()
			if err != nil {
				return model.Task{}, err
			}
			deputies = append(deputies, mappedDeputy)
		}
	}

	return model.Task{
		Assignee:      r.Assignee.toAssignee(),
		Orders:        orders,
		Clients:       clients,
		Deputies:      deputies,
		DueDate:       r.DueDate,
		Id:            r.Id,
		Type:          r.Type,
		Name:          r.Name,
		CaseOwnerTask: r.CaseOwnerTask,
		IsPriority:    r.IsPriority,
	}, nil
}

type taskListResponse struct {
	Tasks      []taskResponse          `json:"tasks"`
	Pages      pageInformationResponse `json:"pages"`
	TotalTasks int                     `json:"total"`
	MetaData   taskMetaDataResponse    `json:"metadata"`
}

func (r taskListResponse) toTaskList() (TaskList, error) {
	var tasks []model.Task
	if r.Tasks != nil {
		tasks = make([]model.Task, 0, len(r.Tasks))
		for _, task := range r.Tasks {
			mappedTask, err := task.toTask()
			if err != nil {
				return TaskList{}, err
			}
			tasks = append(tasks, mappedTask)
		}
	}

	return TaskList{
		Tasks:      tasks,
		Pages:      r.Pages.toPageInformation(),
		TotalTasks: r.TotalTasks,
		MetaData: TaskMetaData{
			TaskTypeCount: r.MetaData.TaskTypeCount,
			AssigneeCount: r.MetaData.AssigneeCount,
		},
	}, nil
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
		"/v1/assignees/teams/tasks?%s&filter=%s&limit=%d&page=%d&sort=%s",
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

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		c.logResponse(req, resp, err)
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logResponse(req, resp, err)
		return v, newStatusError(resp)
	}

	var response taskListResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logResponse(req, resp, err)
		return v, err
	}

	v, err = response.toTaskList()
	if err != nil {
		c.logResponse(req, resp, err)
		return TaskList{}, err
	}

	return v, nil
}

func (p TaskListParams) CreateFilter() string {
	filter := "status:Not+started,"

	if slices.Contains(p.SelectedTaskTypes, TaskTypeEcmHandle) {
		p.SelectedTaskTypes = getEcmTaskTypesString(p.TaskTypes)
	}
	for _, t := range p.SelectedTaskTypes {
		filter += "type:" + t + ","
	}
	for _, a := range p.Assignees {
		filter += "assigneeid_or_null:" + a + ","
	}
	for _, d := range p.Deputies {
		filter += "deputyid:" + d + ","
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
