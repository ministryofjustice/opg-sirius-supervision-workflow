package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ClosedClientsParams struct {
	TeamIds []string `json:"teamIds"`
}

func (c *ApiClient) GetClosedClientList(ctx Context, params ClientListParams) (ClientList, error) {
	var v ClientList
	var filter string
	var body bytes.Buffer
	var err error

	filter = params.CreateFilter()
	ClosedClientMemberIds := ClosedClientsParams{TeamIds: CreateMemberIdArray(params)}

	err = json.NewEncoder(&body).Encode(ClosedClientMemberIds)
	if err != nil {
		return v, err
	}

	endpoint := fmt.Sprintf(
		"/api/v1/assignees/closed-clients?limit=%d&page=%d&filter=%s",
		params.PerPage,
		params.Page,
		filter,
	)

	req, err := c.newRequest(ctx, http.MethodGet, endpoint, &body)

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

func CreateMemberIdArray(params ClientListParams) []string {
	var teamMemberIds []string
	teamMemberIds = append(teamMemberIds, strconv.Itoa(params.Team.Id))
	for _, member := range params.Team.Members {
		teamMemberIds = append(teamMemberIds, strconv.Itoa(member.Id))
	}
	return teamMemberIds
}
