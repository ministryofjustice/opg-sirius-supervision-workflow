package sirius

import (
	"encoding/json"
	"net/http"
)

type PermissionGroup struct {
	Permissions []string `json:"permissions"`
}

type PermissionSet map[string]PermissionGroup

func (c *Client) MyPermissions(ctx Context) (PermissionSet, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/permissions", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v PermissionSet
	err = json.NewDecoder(resp.Body).Decode(&v)

	return PermissionSet{}, err
}
