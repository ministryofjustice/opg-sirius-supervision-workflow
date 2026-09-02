package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (c *ApiClient) GetClosedClientList(ctx Context, params ClientListParams) (ClientList, error) {
	var v ClientList
	query := url.Values{}

	for _, teamId := range CreateMemberIdArray(params) {
		query.Add("teamIds[]", teamId)
	}
	query.Set("limit", strconv.Itoa(params.PerPage))
	query.Set("page", strconv.Itoa(params.Page))
	query.Set("filter", params.CreateFilter())

	endpoint := fmt.Sprintf(
		"/v1/assignees/closed-clients?%s",
		query.Encode(),
	)

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
