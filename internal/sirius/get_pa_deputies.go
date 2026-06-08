package sirius

import (
	"encoding/json"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
)

// GetPADeputies returns every deputy whose executive case manager is in
// any team of type PA. The endpoint returns a slim id/displayName payload
// (a bare JSON array) so the response decodes into model.Deputy values
// with only Id and DisplayName populated; other Deputy fields are left
// at their zero values.
func (c *ApiClient) GetPADeputies(ctx Context) ([]model.Deputy, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/assignees/pa-deputies", nil)
	if err != nil {
		c.logErrorRequest(req, err)
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.logResponse(req, resp, err)
		return nil, err
	}
	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		c.logResponse(req, resp, err)
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logResponse(req, resp, err)
		return nil, newStatusError(resp)
	}

	var deputies []model.Deputy
	if err = json.NewDecoder(resp.Body).Decode(&deputies); err != nil {
		c.logResponse(req, resp, err)
		return nil, err
	}

	return deputies, nil
}
