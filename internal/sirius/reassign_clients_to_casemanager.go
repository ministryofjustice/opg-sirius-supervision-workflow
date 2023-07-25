package sirius

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ReassignClientDetails struct {
	AssigneeId int      `json:"assigneeId"`
	ClientIds  []string `json:"clientIds"`
	IsWorkflow bool     `json:"isWorkflow"`
}

type ReAssigneeResponse struct {
	ReAssigneeName string `json:"reAssigneeName"`
}

func (c *ApiClient) ReassignClientToCaseManager(ctx Context, newAssigneeIdForClient int, clientIds []string) (string, error) {
	var u ReAssigneeResponse
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(ReassignClientDetails{
		AssigneeId: newAssigneeIdForClient,
		ClientIds:  clientIds,
		IsWorkflow: true,
	})

	if err != nil {
		return "", err
	}
	req, err := c.newRequest(ctx, http.MethodPut, "/api/v1/clients/edit/reassign", &body)

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
	fmt.Println(u)
	return u.ReAssigneeName, err
}
