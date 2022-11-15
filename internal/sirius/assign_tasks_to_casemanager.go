package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) AssignTasksToCaseManager(ctx Context, newAssigneeIdForTask int, taskIdForUrl string) error {

	requestURL := fmt.Sprintf("/api/v1/users/%d/tasks/%s", newAssigneeIdForTask, taskIdForUrl)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, nil)

	if err != nil {
		c.logErrorRequest(req, err)
		return err
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		c.logResponse(req, resp, err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.logResponse(req, resp, err)
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			c.logResponse(req, resp, err)
			return &ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}
	return nil
}
