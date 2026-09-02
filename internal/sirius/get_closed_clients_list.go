package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (c *ApiClient) GetClosedClientList(ctx Context, params ClientListParams) (ClientList, error) {
	var v ClientList
	var teamIds []string

	for _, teamId := range CreateMemberIdArray(params) {
		teamIds = append(teamIds, "teamIds[]="+teamId)
	}

	endpoint := fmt.Sprintf(
		"/v1/assignees/closed-clients?%s&limit=%d&page=%d&filter=%s",
		strings.Join(teamIds, "&"),
		params.PerPage,
		params.Page,
		params.CreateFilter(),
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
