package sirius

import (
	"encoding/json"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
)

type currentUserResponse struct {
	ID          int            `json:"id"`
	Name        string         `json:"name,omitempty"`
	PhoneNumber string         `json:"phoneNumber"`
	Teams       []teamResponse `json:"teams"`
	DisplayName string         `json:"displayName"`
	Deleted     bool           `json:"deleted"`
	Email       string         `json:"email"`
	Firstname   string         `json:"firstname"`
	Surname     string         `json:"surname"`
	Roles       []string       `json:"roles"`
	Locked      bool           `json:"locked"`
	Suspended   bool           `json:"suspended"`
}

func (r currentUserResponse) toAssignee() model.Assignee {
	var teams []model.Team
	if r.Teams != nil {
		teams = make([]model.Team, 0, len(r.Teams))
		for _, team := range r.Teams {
			teams = append(teams, team.toTeam())
		}
	}

	return model.Assignee{
		Id:          r.ID,
		Name:        r.DisplayName,
		PhoneNumber: r.PhoneNumber,
		Teams:       teams,
		Deleted:     r.Deleted,
		Email:       r.Email,
		Firstname:   r.Firstname,
		Surname:     r.Surname,
		Roles:       r.Roles,
		Locked:      r.Locked,
		Suspended:   r.Suspended,
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

	return response.toAssignee(), nil
}
