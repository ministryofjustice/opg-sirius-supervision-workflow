package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
)

type ClientListParams struct {
	Team              model.Team
	Page              int
	PerPage           int
	CaseOwners        []string
	OrderStatuses     []string
	SubType           string
	DeputyTypes       []string
	CaseTypes         []string
	SupervisionLevels []string
}

type ClientMetaData struct {
	AssigneeCount []model.AssigneeAndCount `json:"assigneeClientCount"`
}

type ClientList struct {
	Clients      []model.Client        `json:"clients"`
	Pages        model.PageInformation `json:"pages"`
	TotalClients int                   `json:"total"`
	MetaData     ClientMetaData        `json:"metadata"`
}

type clientListResponse struct {
	Clients      []clientResponse        `json:"clients"`
	Pages        pageInformationResponse `json:"pages"`
	TotalClients int                     `json:"total"`
	MetaData     ClientMetaData          `json:"metadata"`
}

func (r clientListResponse) toClientList() (ClientList, error) {
	var clients []model.Client
	if r.Clients != nil {
		clients = make([]model.Client, 0, len(r.Clients))
		for _, client := range r.Clients {
			mappedClient, err := client.toClient()
			if err != nil {
				return ClientList{}, err
			}
			clients = append(clients, mappedClient)
		}
	}

	return ClientList{
		Clients:      clients,
		Pages:        r.Pages.toPageInformation(),
		TotalClients: r.TotalClients,
		MetaData:     r.MetaData,
	}, nil
}

func (c *ApiClient) GetClientList(ctx Context, params ClientListParams) (ClientList, error) {
	var v ClientList
	var sort string
	var filter string

	if params.Team.IsLay() {
		sort = "report_due_date:asc"
	}

	if params.Team.IsLayNewOrdersTeam() {
		sort = "made_active_date:asc"
	} else {
		filter = params.CreateFilter()
	}

	endpoint := fmt.Sprintf("/v1/assignees/%d/clients?limit=%d&page=%d&filter=%s&sort=%s", params.Team.Id, params.PerPage, params.Page, filter, sort)
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

	var response clientListResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logResponse(req, resp, err)
		return v, err
	}

	v, err = response.toClientList()
	if err != nil {
		c.logResponse(req, resp, err)
		return ClientList{}, err
	}

	return v, err
}

func (p ClientListParams) CreateFilter() string {
	var filter string
	for _, s := range p.OrderStatuses {
		filter += "order-status:" + s + ","
	}
	if p.SubType != "" {
		filter += "subtype:" + p.SubType + ","
	}
	for _, dt := range p.DeputyTypes {
		filter += "deputy-type:" + dt + ","
	}
	for _, ct := range p.CaseTypes {
		filter += "case-type:" + ct + ","
	}
	for _, a := range p.CaseOwners {
		filter += "caseowner:" + a + ","
	}
	for _, a := range p.SupervisionLevels {
		filter += "supervision-level:" + a + ","
	}
	return strings.TrimRight(filter, ",")
}
