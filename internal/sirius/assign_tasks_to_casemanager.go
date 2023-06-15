package sirius

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type ReassignTaskDetails struct {
	AssigneeId int      `json:"assigneeId"`
	TaskIds    []string `json:"taskIds"`
	IsPriority string   `json:"isPriority"`
}

func (c *ApiClient) AssignTasksToCaseManager(ctx Context, newAssigneeIdForTask int, taskIds []string, prioritySelected string) (string, error) {
	var u ApiTask
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(ReassignTaskDetails{
		AssigneeId: newAssigneeIdForTask,
		TaskIds:    taskIds,
		IsPriority: prioritySelected,
	})

	if err != nil {
		return "", err
	}
	req, err := c.newRequest(ctx, http.MethodPut, "/api/v1/reassign-tasks", &body)

	if err != nil {
		c.logErrorRequest(req, err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		c.logResponse(req, resp, err)
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", ErrUnauthorized
	}

	if resp.StatusCode == http.StatusForbidden {
		return "", errors.New("Only managers can set priority on tasks")
	}

	if resp.StatusCode != http.StatusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			c.logResponse(req, resp, err)
			return "", &ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return "", newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&u)
	if err != nil {
		c.logResponse(req, resp, err)
		return "", err
	}
	return u.ApiTaskAssignee.CaseManagerName, err
}
