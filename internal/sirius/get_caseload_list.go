package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AnnualReport struct {
	DueDate string `json:"dueDate"`
}

type Order struct {
	OrderId            int          `json:"id"`
	OrderStatus        RefData      `json:"orderStatus"`
	LatestAnnualReport AnnualReport `json:"latestAnnualReport"`
}

type ApiClient struct {
	ClientId                   int            `json:"id"`
	ClientCaseRecNumber        string         `json:"caseRecNumber"`
	ClientFirstName            string         `json:"firstname"`
	ClientSurname              string         `json:"surname"`
	ClientSupervisionCaseOwner CaseManagement `json:"supervisionCaseOwner"`
	Case                       []Order        `json:"cases"`
	SupervisionLevel           string         `json:"supervisionLevel"`
}

type ClientList struct {
	WholeClientList []ApiClient     `json:"clients"`
	Pages           PageInformation `json:"pages"`
	TotalClients    int             `json:"total"`
}

func (c *Client) GetCaseloadList(ctx Context, teamIds string) (ClientList, error) {
	var v ClientList

	endpoint := fmt.Sprintf("/api/v1/assignees/%s/clients", teamIds)
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
