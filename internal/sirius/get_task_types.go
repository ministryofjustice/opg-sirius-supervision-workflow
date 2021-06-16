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
	IsSelected bool
}

type WholeTaskList struct {
	AllTaskList map[string]ApiTaskTypes `json:"task_types"`
}

func (c *Client) GetTaskTypes(ctx Context, taskTypeSelected []string) ([]ApiTaskTypes, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/tasktypes/supervision", nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v WholeTaskList
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	WholeTaskList := v.AllTaskList

	var taskTypeList []ApiTaskTypes

	for _, u := range WholeTaskList {
		taskType := ApiTaskTypes{
			Handle:     u.Handle,
			Incomplete: u.Incomplete,
			Category:   u.Category,
			Complete:   u.Complete,
			User:       u.User,
			IsSelected: IsSelected(u.Handle, taskTypeSelected),
		}
		taskTypeList = append(taskTypeList, taskType)
	}

	sort.Slice(taskTypeList, func(i, j int) bool {
		return taskTypeList[i].Incomplete < taskTypeList[j].Incomplete
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
