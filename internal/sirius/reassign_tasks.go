package sirius

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"net/http"
	"strconv"
)

type ReassignTasksParams struct {
	AssignTeam string
	AssignCM   string
	AssigneeId int      `json:"assigneeId"`
	TaskIds    []string `json:"taskIds"`
	IsPriority string   `json:"isPriority"`
}

func (c *ApiClient) ReassignTasks(ctx Context, params ReassignTasksParams) (string, error) {
	var u model.Task
	var body bytes.Buffer
	var err error

	assignee := params.AssignTeam
	if params.AssignCM != "" {
		assignee = params.AssignCM
	}

	params.AssigneeId, err = strconv.Atoi(assignee)
	if err != nil {
		return "", err
	}

	err = json.NewEncoder(&body).Encode(params)

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

	if params.AssignTeam != "0" {
		switch params.IsPriority {
		case "true":
			return fmt.Sprintf("You have assigned %d task(s) to %s as a priority", len(params.TaskIds), u.Assignee.Name), nil
		case "false":
			return fmt.Sprintf("You have assigned %d task(s) to %s and removed priority", len(params.TaskIds), u.Assignee.Name), nil
		default:
			return fmt.Sprintf("You have assigned %d task(s) to %s", len(params.TaskIds), u.Assignee.Name), nil
		}
	}
	switch params.IsPriority {
	case "true":
		return fmt.Sprintf("You have assigned %d task(s) as a priority", len(params.TaskIds)), nil
	case "false":
		return fmt.Sprintf("You have removed %d task(s) as a priority", len(params.TaskIds)), nil
	}

	return "", nil
}
