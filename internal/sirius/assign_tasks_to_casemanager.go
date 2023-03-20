package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ReassignTaskDetails struct {
	UserId  int      `json:"userId"`
	TaskIds []string `json:"taskIds"`
}

func (c *Client) AssignTasksToCaseManager(ctx Context, newAssigneeIdForTask int, taskIds []string) error {

	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(ReassignTaskDetails{
		UserId:  newAssigneeIdForTask,
		TaskIds: taskIds,
	})

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/reassign-multiple-tasks"), &body)

	if err != nil {
		c.logErrorRequest(req, err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		c.logResponse(req, resp, err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			c.logResponse(req, resp, err)
			return &ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}
