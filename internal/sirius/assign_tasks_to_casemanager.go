package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-go-common/logging"
	"net/http"
)

func (c *Client) AssignTasksToCaseManager(ctx Context, logger *logging.Logger, newAssigneeIdForTask int, taskIdForUrl string) error {

	requestURL := fmt.Sprintf("/api/v1/users/%d/tasks/%s", newAssigneeIdForTask, taskIdForUrl)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, nil)
	c.logRequest(req, err)

	if err != nil {
		return err
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
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

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return &ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}
	return nil
}
