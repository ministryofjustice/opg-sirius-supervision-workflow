package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"net/http"
	"strings"
	"time"
)

type ClientListParams struct {
	Team             model.Team
	Page             int
	PerPage          int
	CaseOwners       []string
	OrderStatuses    []string
	SubType          string
	DebtTypes        []string
	DeputyTypes      []string
	CaseTypes        []string
	LastActionDate   string
	CachedDebtAmount string
}

type ClientMetadata []struct {
	ClientId         int    `json:"clientId"`
	LastActionDate   string `json:"lastActionDate"`
	CachedDebtAmount string `json:"cachedDebtAmount"`
}

type ClientList struct {
	Clients            []model.Client        `json:"clients"`
	Pages              model.PageInformation `json:"pages"`
	TotalClients       int                   `json:"total"`
	ClientListMetaData ClientMetadata        `json:"metadata"`
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
	if params.Team.IsClosedCases() {
		endpoint = fmt.Sprintf("/api/v1/assignees/%d/closed-clients?limit=%d&page=%d&filter=%s&sort=%s", params.Team.Id, params.PerPage, params.Page, filter, sort)
	}

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
	//io.Copy(os.Stdout, resp.Body)

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

	v = appendMetaData(v)
	return v, err
}

func appendMetaData(v ClientList) ClientList {
	for i, s := range v.Clients {
		for _, t := range v.ClientListMetaData {
			if s.Id == t.ClientId {
				stringDate := formatTimestampToStandardDate(t.LastActionDate)
				v.Clients[i].LastActionDate = stringDate
				v.Clients[i].CachedDebtAmount = t.CachedDebtAmount
			}
		}
	}
	fmt.Println("append meta data")
	fmt.Println(v.Clients[0].LastActionDate)
	fmt.Println(v.Clients[0].CachedDebtAmount)

	return v
}

func formatTimestampToStandardDate(timestamp string) string {
	newTime, _ := time.Parse("2006-01-02T15:04:05Z07:00", timestamp)
	return newTime.Format("02/01/2006")
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
	for _, d := range p.DebtTypes {
		filter += "debt:" + d + ","
	}
	return strings.TrimRight(filter, ",")
}
