package sirius

import (
	"encoding/json"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"

	"net/http"
)

func (c *ApiClient) GetCurrentUserDetails(ctx Context) (model.Assignee, error) {
	var user model.Assignee

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/users/current", nil)
	if err != nil {
		c.logErrorRequest(req, err)
		return user, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.logRequest(req, err)
		return user, err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		c.logRequest(req, err)
		return user, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logRequest(req, err)
		return user, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	return user, err
}
