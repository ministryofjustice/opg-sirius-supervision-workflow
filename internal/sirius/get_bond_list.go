package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
)

type BondMetaData struct {
	BondMetaData []model.AssigneeAndCount `json:"ecmCount"`
}

type BondList struct {
	Bonds      []model.Bond          `json:"bonds"`
	Pages      model.PageInformation `json:"pages"`
	TotalBonds int                   `json:"total"`
}

type BondListParams struct {
	Team    model.Team
	Page    int
	PerPage int
}

func (c *ApiClient) GetBondList(ctx Context, params BondListParams) (BondList, error) {
	var v BondList

	endpoint := fmt.Sprintf("/v1/bonds/without-orders?limit=%d&page=%d", params.PerPage, params.Page)
	req, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)

	if err != nil {
		c.logErrorRequest(req, err)
		return v, err
	}

	resp, err := c.http.Do(req)

	if err != nil {
		c.logResponse(req, resp, err)
		return v, err
	}
	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		c.logResponse(req, resp, err)
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logResponse(req, resp, err)
		return v, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		c.logResponse(req, resp, err)
		return v, err
	}

	return v, nil
}
