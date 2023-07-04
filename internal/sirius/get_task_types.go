package sirius

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"net/http"
	"sort"
)

type TaskTypesList struct {
	TaskTypes map[string]model.TaskType `json:"task_types"`
}

func (c *ApiClient) GetTaskTypes(ctx Context, selectedTaskTypes []string) ([]model.TaskType, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/tasktypes/supervision", nil)

	if err != nil {
		c.logErrorRequest(req, err)
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.logResponse(req, resp, err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.logResponse(req, resp, err)
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logResponse(req, resp, err)
		return nil, newStatusError(resp)
	}

	var v TaskTypesList
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		c.logResponse(req, resp, err)
		return nil, err
	}

	var taskTypes []model.TaskType

	for _, u := range v.TaskTypes {
		taskType := model.TaskType{
			Handle:     u.Handle,
			Incomplete: u.Incomplete,
			Category:   u.Category,
			Complete:   u.Complete,
			User:       u.User,
			EcmTask:    u.EcmTask,
			IsSelected: IsSelected(u.Handle, selectedTaskTypes),
		}
		taskTypes = append(taskTypes, taskType)
	}

	sort.Slice(taskTypes, func(i, j int) bool {
		return taskTypes[i].Incomplete < taskTypes[j].Incomplete
	})

	// prepend the "ECM Tasks" filter option
	taskTypes = append([]model.TaskType{
		{
			Handle:     "ECM_TASKS",
			Incomplete: "ECM Tasks",
			IsSelected: IsSelected("ECM_TASKS", selectedTaskTypes),
		},
	}, taskTypes...)

	return taskTypes, err
}

func IsSelected(handle string, taskTypeSelected []string) bool {
	for _, q := range taskTypeSelected {
		if handle == q {
			return true
		}
	}
	return false
}
