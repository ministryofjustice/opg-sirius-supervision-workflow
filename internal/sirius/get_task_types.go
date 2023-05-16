package sirius

import (
	"encoding/json"
	"net/http"
	"sort"
)

type TaskTypes struct {
	Handle     string `json:"handle"`
	Name       string `json:"incomplete"`
	IsSelected bool
	TaskCount  int
}

type WholeTaskTypesList struct {
	AllTaskList map[string]TaskTypes `json:"task_types"`
}

func (c *Client) GetTaskTypes(ctx Context, taskTypeSelected []string) ([]TaskTypes, error) {
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

	var taskTypeList []TaskTypes

	for _, u := range v.AllTaskList {
		taskType := TaskTypes{
			Handle:     u.Handle,
			Name:       u.Name,
			IsSelected: IsSelected(u.Handle, taskTypeSelected),
		}
		taskTypeList = append(taskTypeList, taskType)
	}

	sort.Slice(taskTypeList, func(i, j int) bool {
		return taskTypeList[i].Name < taskTypeList[j].Name
	})

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
