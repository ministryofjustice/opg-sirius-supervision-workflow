package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) DeleteUser(ctx Context, userID int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/auth/user/%d", userID), nil)
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

	if resp.StatusCode != http.StatusOK {
		var v struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ClientError(v.Message)
		}

		return newStatusError(resp)
	}

	return nil
}
