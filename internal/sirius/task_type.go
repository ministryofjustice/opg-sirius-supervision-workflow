package sirius

import (
	"encoding/json"
	"net/http"
)

type ApiTaskTypes struct {
	Category   string `json:"category"`
	Complete   string `json:"complete"`
	Handle     string `json:"handle"`
	Incomplete string `json:"incomplete"`
	User       bool   `json:"user"`
}

type LoadTasks struct {
	Category   string
	Complete   string
	Handle     string
	Incomplete string
	User       bool
}

func (c *Client) GetTaskDetails(ctx Context) ([]LoadTasks, error) {

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

	var v []ApiTaskTypes
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	tasklist := make([]LoadTasks, len(v))

	// if resp.StatusCode == http.StatusUnauthorized {
	// 	return nil, ErrUnauthorized
	// }

	// if resp.StatusCode != http.StatusOK {
	// 	return nil, newStatusError(resp)
	// }

	// var v []apiTaskTypes
	// if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
	// 	return nil, err
	// }

	// taskTypeList := make([]LoadTasks, len(v))

	// for i, t := range v {
	// 	taskTypeList[i] = LoadTasks{
	// 		Category:   t.Category,
	// 		Complete:   t.Complete,
	// 		Handle:     t.Handle,
	// 		Incomplete: t.Incomplete,
	// 		User:       t.User,
	// 	}
	// }

	return tasklist, nil
}
