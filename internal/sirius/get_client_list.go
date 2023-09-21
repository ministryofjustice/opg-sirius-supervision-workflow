package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"net/http"
	"strings"
)

type ClientListParams struct {
	Team          model.Team
	Page          int
	PerPage       int
	CaseOwners    []string
	OrderStatuses []string
	SubType       string
}

type ClientList struct {
	Clients      []model.Client        `json:"clients"`
	Pages        model.PageInformation `json:"pages"`
	TotalClients int                   `json:"total"`
}

func (c *ApiClient) GetClientList(ctx Context, params ClientListParams) (ClientList, error) {
	var v ClientList
	var sort string
	var filter string

	if params.Team.IsLayNewOrdersTeam() {
		sort = "made_active_date:asc"
	} else {
		filter = params.CreateFilter()
	}

	endpoint := fmt.Sprintf("/api/v1/assignees/%d/clients?limit=%d&page=%d&filter=%s&sort=%s", params.Team.Id, params.PerPage, params.Page, filter, sort)
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

func (p ClientListParams) CreateFilter() string {
	var filter string
	for _, a := range p.CaseOwners {
		filter += "caseowner:" + a + ","
	}
	for _, s := range p.OrderStatuses {
		filter += "order-status:" + s + ","
	}
	return strings.TrimRight(filter, ",")
}
