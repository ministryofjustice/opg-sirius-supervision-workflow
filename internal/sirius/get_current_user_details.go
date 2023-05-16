package sirius

import (
	"encoding/json"

	"net/http"
)

type UserDetails struct {
	Firstname string      `json:"firstname"`
	Surname   string      `json:"surname"`
	Teams     []UserTeams `json:"teams"`
}

type UserTeams struct {
	TeamId int `json:"id"`
}

func (c *Client) GetCurrentUserDetails(ctx Context) (UserDetails, error) {
	var v UserDetails

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/users/current", nil)
	if err != nil {
		c.logErrorRequest(req, err)
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.logger.Request(req, err)
		return v, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.logger.Request(req, err)
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Request(req, err)
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}
