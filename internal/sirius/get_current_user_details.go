package sirius

import (
	"encoding/json"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
)

type currentUserResponse struct {
	ID          int            `json:"id"`
	DisplayName string         `json:"displayName"`
	Teams       []teamResponse `json:"teams"`
	Roles       []string       `json:"roles"`
}

func (r currentUserResponse) model() model.Assignee {
	var teams []model.Team
	if r.Teams != nil {
		teams = make([]model.Team, 0, len(r.Teams))
		for _, team := range r.Teams {
			teams = append(teams, team.model())
		}
	}

	return model.Assignee{
		Id:    r.ID,
		Name:  r.DisplayName,
		Teams: teams,
		Roles: r.Roles,
	}
}

func (c *ApiClient) GetCurrentUserDetails(ctx Context) (model.Assignee, error) {
	var user model.Assignee

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/users/current", nil)
	if err != nil {
		c.logErrorRequest(req, err)
		return user, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.logRequest(req, err)
		return user, err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		c.logRequest(req, err)
		return user, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logRequest(req, err)
		return user, newStatusError(resp)
	}

	var response currentUserResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return user, err
	}

	return response.model(), nil
}
