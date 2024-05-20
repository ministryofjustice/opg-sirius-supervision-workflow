package sirius

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"

	"net/http"
)

func (c *ApiClient) GetCurrentUserDetails(ctx Context) (model.Assignee, error) {
	var v model.Assignee

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/users/current", nil)
	if err != nil {
		c.logErrorRequest(req, err)
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.logRequest(req, err)
		return v, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.logRequest(req, err)
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logRequest(req, err)
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}
