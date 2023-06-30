package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"net/http"
)

type ClientList struct {
	Clients      []model.Client        `json:"clients"`
	Pages        model.PageInformation `json:"pages"`
	TotalClients int                   `json:"total"`
}

func (c *ApiClient) GetClientList(ctx Context, team model.Team) (ClientList, error) {
	var v ClientList

	var sort string
	if team.IsLayNewOrdersTeam() {
		sort = "made_active_date:asc"
	}

	endpoint := fmt.Sprintf("/api/v1/assignees/%d/clients?sort=%s", team.Id, sort)
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

	defer resp.Body.Close()

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

	return v, err
}
