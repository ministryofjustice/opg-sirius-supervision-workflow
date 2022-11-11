package sirius

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-go-common/logging"

	"net/http"
)

type UserDetails struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	PhoneNumber string          `json:"phoneNumber"`
	Teams       []MyDetailsTeam `json:"teams"`
	DisplayName string          `json:"displayName"`
	Deleted     bool            `json:"deleted"`
	Email       string          `json:"email"`
	Firstname   string          `json:"firstname"`
	Surname     string          `json:"surname"`
	Roles       []string        `json:"roles"`
	Locked      bool            `json:"locked"`
	Suspended   bool            `json:"suspended"`
}

type MyDetailsTeam struct {
	DisplayName string `json:"displayName"`
	TeamId      int    `json:"id"`
}

func (c *Client) GetCurrentUserDetails(ctx Context, logger *logging.Logger) (UserDetails, error) {
	var v UserDetails

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/users/current", nil)
	c.logRequest(logger, req, err)

	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}
