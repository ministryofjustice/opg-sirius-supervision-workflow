package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type deputyListResponse struct {
	Deputies      []deputyResponse        `json:"persons"`
	Pages         pageInformationResponse `json:"pages"`
	TotalDeputies int                     `json:"total"`
	MetaData      DeputyMetaData          `json:"metadata"`
}

func (r deputyListResponse) toDeputyList() (DeputyList, error) {
	var deputies []model.Deputy
	if r.Deputies != nil {
		deputies = make([]model.Deputy, 0, len(r.Deputies))
		for _, deputy := range r.Deputies {
			mappedDeputy, err := deputy.toDeputy()
			if err != nil {
				return DeputyList{}, err
			}
			deputies = append(deputies, mappedDeputy)
		}
	}

	return DeputyList{
		Deputies:      deputies,
		Pages:         r.Pages.toPageInformation(),
		TotalDeputies: r.TotalDeputies,
		MetaData:      r.MetaData,
	}, nil
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
		"/v1/assignees/teams/deputies?%s&limit=%d&page=%d&filter=%s&sort=%s",
		strings.Join(teamIds, "&"),
		params.PerPage,
		params.Page,
		params.CreateFilter(),
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

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		c.logResponse(req, resp, err)
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logResponse(req, resp, err)
		return v, newStatusError(resp)
	}

	var response deputyListResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logResponse(req, resp, err)
		return v, err
	}

	v, err = response.toDeputyList()
	if err != nil {
		c.logResponse(req, resp, err)
		return DeputyList{}, err
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
