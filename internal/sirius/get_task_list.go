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

func (c *ApiClient) GetTaskList(ctx Context, search int, displayTaskLimit int, selectedTeam model.Team, taskTypeSelected []string, taskTypes []model.TaskType, selectedAssignees []string, dueDateFrom *time.Time, dueDateTo *time.Time) (TaskList, error) {
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

func CreateFilter(taskTypeSelected []string, selectedAssignees []string, dueDateFrom *time.Time, dueDateTo *time.Time, taskTypes []model.TaskType) string {
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

func getEcmTaskTypesString(taskTypes []model.TaskType) []string {
	var ecmTasks []string
	for _, taskType := range taskTypes {
		if taskType.EcmTask {
			ecmTasks = append(ecmTasks, taskType.Handle)
		}
	}
	return ecmTasks
}
