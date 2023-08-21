package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"net/http"
	"sort"
)

const TaskTypeCategorySupervision = "supervision"
const TaskTypeCategoryDeputy = "deputy"
const TaskTypeEcmHandle = "ECM_TASKS"
const TaskTypeEcmLabel = "ECM Tasks"

type TaskTypesList struct {
	TaskTypes map[string]model.TaskType `json:"task_types"`
}

type TaskTypesParams struct {
	Category  string
	ProDeputy bool
	PADeputy  bool
}

func (c *ApiClient) GetTaskTypes(ctx Context, params TaskTypesParams) ([]model.TaskType, error) {
	endpoint := fmt.Sprintf("/api/v1/tasktypes/%s", params.Category)
	if params.ProDeputy {
		endpoint += "?pro_deputy=true"
	} else if params.PADeputy {
		endpoint += "?pa_deputy=true"
	}

	req, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)

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
		taskTypes = append(taskTypes, u)
	}

	sort.Slice(taskTypes, func(i, j int) bool {
		return taskTypes[i].Incomplete < taskTypes[j].Incomplete
	})

	if params.Category == TaskTypeCategorySupervision {
		taskTypes = append([]model.TaskType{
			{
				Handle:     TaskTypeEcmHandle,
				Incomplete: TaskTypeEcmLabel,
			},
		}, taskTypes...)
	}

	return taskTypes, nil
}
