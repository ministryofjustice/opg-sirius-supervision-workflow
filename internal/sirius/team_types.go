package sirius

import (
	"encoding/json"
	"net/http"
)

type RefDataTeamType struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (c *Client) TeamTypes(ctx Context) ([]RefDataTeamType, error) {
	var v struct {
		Data []RefDataTeamType `json:"teamType"`
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/reference-data?filter=teamType", nil)
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
	return v.Data, err
}
