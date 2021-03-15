package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type addTeamRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

func (c *Client) AddTeam(ctx Context, name, teamType, phone, email string) (int, error) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(addTeamRequest{
		Name:        name,
		Email:       email,
		PhoneNumber: phone,
		Type:        teamType,
	})
	if err != nil {
		return 0, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/teams", &body)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return 0, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusCreated {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return 0, ValidationError{Errors: v.ValidationErrors}
		}

		return 0, newStatusError(resp)
	}

	var v apiTeam
	err = json.NewDecoder(resp.Body).Decode(&v)

	return v.ID, err
}
