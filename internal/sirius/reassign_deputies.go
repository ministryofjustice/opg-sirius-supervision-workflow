package sirius

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type ReassignDeputiesParams struct {
	AssignTeam string
	AssignCM   string
	AssigneeId int      `json:"assigneeId"`
	DeputyIds  []string `json:"deputyIds"`
	IsWorkflow bool     `json:"isWorkflow"`
}

func (c *ApiClient) ReassignDeputies(ctx Context, params ReassignDeputiesParams) (string, error) {
	var u ReassignResponse
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

	params.IsWorkflow = true
	err = json.NewEncoder(&body).Encode(params)

	if err != nil {
		return "", err
	}
	req, err := c.newRequest(ctx, http.MethodPut, "/api/v1/deputies/reassign", &body)

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
		return "", errors.New("Only managers can reassign deputy cases")
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

	return fmt.Sprintf("You have reassigned %d deputies(s) to %s", len(params.DeputyIds), u.ReassignName), nil
}
