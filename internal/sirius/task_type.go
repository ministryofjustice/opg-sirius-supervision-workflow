package sirius

import (
	"encoding/json"
	"net/http"
)

// type Salary struct {
// 	Basic, HRA, TA float64
// }

// type Employee struct {
// 	FirstName, LastName, Email string
// 	Age                        int
// 	MonthlySalary              []Salary
// }

type ApiTaskTypes struct {
	Handle     string `json:"handle"`
	Incomplete string `json:"incomplete"`
	Category   string `json:"category"`
	Complete   string `json:"complete"`
	User       bool   `json:"user"`
}

type WholeTaskList struct {
	AllTaskList map[string]ApiTaskTypes `json:"task_types"`
}

func (c *Client) GetTaskDetails(ctx Context) (WholeTaskList, error) {
	var v WholeTaskList

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/tasktypes/supervision", nil)
	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	// io.Copy(os.Stdout, resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, err
	}

	// tasklist := make([]LoadTasks, len(v))

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

	return v, err
}
