package sirius

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type ReassignClientDetails struct {
	AssigneeId int      `json:"assigneeId"`
	ClientIds  []string `json:"clientIds"`
	IsWorkflow bool     `json:"isWorkflow"`
}

type ReassignResponse struct {
	ReassignName string `json:"reassignName"`
}

func (c *ApiClient) ReassignClientToCaseManager(ctx Context, newAssigneeIdForClient int, clientIds []string) (string, error) {
	var u ReassignResponse
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
		return "", errors.New("Only managers can reassign client cases")
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

	return u.ReassignName, err
}