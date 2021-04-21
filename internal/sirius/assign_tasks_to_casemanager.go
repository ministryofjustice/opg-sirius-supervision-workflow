package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) AssignTasksToCaseManager(ctx Context, newAssigneeIdForTask int, selectedTask int) error {

	var body bytes.Buffer

	log.Println("start of assign function")
	log.Println("newAssigneeIdForTask")
	log.Println(newAssigneeIdForTask)
	log.Println("selectedTask")
	log.Println(selectedTask)

	requestURL := fmt.Sprintf("/api/v1/users/%d/tasks/%d", newAssigneeIdForTask, selectedTask)

	log.Println("request url")
	log.Println(requestURL)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)

	log.Println("body")
	log.Println(&body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	log.Println("resp body")
	log.Println(resp)
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
	log.Println("end of assign task function")
	return nil
}
