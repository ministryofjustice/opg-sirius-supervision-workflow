package sirius

import (
	"encoding/json"
	"net/http"
	"sort"
)

func (c *Client) Roles(ctx Context) ([]string, error) {
	var v struct {
		Data []string `json:"data"`
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/api/role", nil)
	if err != nil {
		return v.Data, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v.Data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v.Data, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v.Data, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	var roles []string
	for _, role := range v.Data {
		if role != "COP User" && role != "OPG User" {
			roles = append(roles, role)
		}
	}

	sort.Strings(roles)

	return roles, err
}
