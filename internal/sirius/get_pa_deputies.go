package sirius

import (
	"encoding/json"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
)

type paDeputyResponse struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

func (r paDeputyResponse) toDeputy() model.Deputy {
	return model.Deputy{
		Id:          r.ID,
		DisplayName: r.DisplayName,
	}
}

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

	var response []paDeputyResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logResponse(req, resp, err)
		return nil, err
	}

	var deputies []model.Deputy
	if response != nil {
		deputies = make([]model.Deputy, 0, len(response))
		for _, deputy := range response {
			deputies = append(deputies, deputy.toDeputy())
		}
	}

	return deputies, nil
}
