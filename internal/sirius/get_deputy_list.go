package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"net/http"
	"strconv"
	"strings"
)

type DeputyList struct {
	Deputies      []model.Deputy        `json:"persons"`
	Pages         model.PageInformation `json:"pages"`
	TotalDeputies int                   `json:"total"`
}

type DeputyListParams struct {
	Team    model.Team
	Page    int
	PerPage int
	Sort    string
}

func (c *ApiClient) GetDeputyList(ctx Context, params DeputyListParams) (DeputyList, error) {
	var v DeputyList
	var teamIds []string

	if params.Team.Id != 0 {
		teamIds = []string{"teamIds[]=" + strconv.Itoa(params.Team.Id)}
	}
	for _, team := range params.Team.Teams {
		teamIds = append(teamIds, "teamIds[]="+strconv.Itoa(team.Id))
	}

	endpoint := fmt.Sprintf(
		"/api/v1/assignees/teams/deputies?%s&limit=%d&page=%d&sort=%s",
		strings.Join(teamIds, "&"),
		params.PerPage,
		params.Page,
		params.Sort,
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

	return v, nil
}