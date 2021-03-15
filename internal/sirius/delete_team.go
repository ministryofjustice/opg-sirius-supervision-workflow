package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) DeleteTeam(ctx Context, teamID int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/teams/%d", teamID), nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusNoContent {
		var v struct {
			Detail string `json:"detail"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ClientError(v.Detail)
		}

		return newStatusError(resp)
	}

	return nil
}
