package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
)

type DeputyMetaData struct {
	DeputyMetaData []model.AssigneeAndCount `json:"ecmCount"`
}

type DeputyList struct {
	Deputies           []model.Deputy        `json:"persons"`
	Pages              model.PageInformation `json:"pages"`
	TotalDeputies      int                   `json:"total"`
	PaProTeamSelection []model.Team
	MetaData           DeputyMetaData `json:"metadata"`
}

type DeputyListParams struct {
	Team         model.Team
	Page         int
	PerPage      int
	Sort         string
	SelectedECMs []string
}

func (c *ApiClient) GetDeputyList(ctx Context, params DeputyListParams) (DeputyList, error) {
	var v DeputyList
	query := url.Values{}

	if params.Team.Id != 0 {
		query.Add("teamIds[]", strconv.Itoa(params.Team.Id))
	}
	for _, team := range params.Team.Teams {
		query.Add("teamIds[]", strconv.Itoa(team.Id))
	}
	query.Set("limit", strconv.Itoa(params.PerPage))
	query.Set("page", strconv.Itoa(params.Page))
	query.Set("filter", params.CreateFilter())
	query.Set("sort", params.Sort)

	endpoint := fmt.Sprintf(
		"/v1/assignees/teams/deputies?%s",
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

	return v, nil
}

func (d DeputyListParams) CreateFilter() string {
	var filter string
	for _, s := range d.SelectedECMs {
		filter += "ecm:" + s + ","
	}
	return strings.TrimRight(filter, ",")
}
