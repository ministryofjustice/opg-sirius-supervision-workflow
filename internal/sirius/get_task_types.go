package sirius

import (
	"encoding/json"
	"net/http"
	"sort"
)

type ApiTaskTypes struct {
	Handle     string `json:"handle"`
	Incomplete string `json:"incomplete"`
	Category   string `json:"category"`
	Complete   string `json:"complete"`
	User       bool   `json:"user"`
	EcmTask    bool   `json:"ecmTask"`
	IsSelected bool
	TaskCount  int
}

type WholeTaskTypesList struct {
	AllTaskList map[string]ApiTaskTypes `json:"task_types"`
}

func (c *ApiClient) GetTaskTypes(ctx Context, taskTypeSelected []string) ([]ApiTaskTypes, error) {
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

	var v WholeTaskTypesList
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		c.logResponse(req, resp, err)
		return nil, err
	}

	WholeTaskTypeList := v.AllTaskList

	var taskTypeList []ApiTaskTypes

	for _, u := range WholeTaskTypeList {
		taskType := ApiTaskTypes{
			Handle:     u.Handle,
			Incomplete: u.Incomplete,
			Category:   u.Category,
			Complete:   u.Complete,
			User:       u.User,
			EcmTask:    u.EcmTask,
			IsSelected: IsSelected(u.Handle, taskTypeSelected),
		}
		taskTypeList = append(taskTypeList, taskType)
	}

	sort.Slice(taskTypeList, func(i, j int) bool {
		return taskTypeList[i].Incomplete < taskTypeList[j].Incomplete
	})

	// prepend the "ECM Tasks" filter option
	taskTypeList = append([]ApiTaskTypes{
		{
			Handle:     "ECM_TASKS",
			Incomplete: "ECM Tasks",
			IsSelected: IsSelected("ECM_TASKS", taskTypeSelected),
		},
	}, taskTypeList...)

	return taskTypeList, err
}

func IsSelected(handle string, taskTypeSelected []string) bool {
	for _, q := range taskTypeSelected {
		if handle == q {
			return true
		}
	}
	return false
}
