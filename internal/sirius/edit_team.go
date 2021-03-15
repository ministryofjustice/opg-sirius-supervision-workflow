package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type editTeamRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	MemberIds   []int  `json:"memberIds"`
}

func (c *Client) EditTeam(ctx Context, team Team) error {
	memberIDs := make([]int, len(team.Members))
	for i, member := range team.Members {
		memberIDs[i] = member.ID
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(editTeamRequest{
		Name:        team.DisplayName,
		Email:       team.Email,
		PhoneNumber: team.PhoneNumber,
		Type:        team.Type,
		MemberIds:   memberIDs,
	})

	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/api/v1/teams/%d", team.ID)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

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
