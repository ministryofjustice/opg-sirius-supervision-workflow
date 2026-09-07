package sirius

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type ReassignClientsParams struct {
	AssignTeam string
	AssignCM   string
	ClientIds  []string
}

type ReassignClientsRequest struct {
	AssigneeId int      `json:"assigneeId"`
	ClientIds  []string `json:"clientIds"`
	IsWorkflow bool     `json:"isWorkflow"`
}

type ReassignResponse struct {
	ReassignName string `json:"reassignName"`
}

func (c *ApiClient) ReassignClients(ctx Context, params ReassignClientsParams) (string, error) {
	var u ReassignResponse
	var body bytes.Buffer
	var err error

	assignee := params.AssignTeam
	if params.AssignCM != "" {
		assignee = params.AssignCM
	}

	id, err := strconv.Atoi(assignee)
	if err != nil {
		return "", err
	}

	err = json.NewEncoder(&body).Encode(ReassignClientsRequest{
		AssigneeId: id,
		ClientIds:  params.ClientIds,
		IsWorkflow: true,
	})

	if err != nil {
		return "", err
	}
	req, err := c.newRequest(ctx, http.MethodPut, "/v1/clients/edit/reassign", &body)

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

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return "", ErrUnauthorized
	}

	if resp.StatusCode == http.StatusForbidden {
		return "", errors.New("only managers can reassign client cases")
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

	return fmt.Sprintf("You have reassigned %d client(s) to %s", len(params.ClientIds), u.ReassignName), nil
}
