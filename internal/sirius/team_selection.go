package sirius

import (
	"encoding/json"
	"net/http"
)

type TeamMembers struct {
	// TeamMembersDeleted      bool   `json:"deleted"`
	// TeamMembersDisplayName  string `json:"displayName"`
	// TeamMembersEmail        string `json:"email"`
	TeamMembersId int `json:"id"`
	// TeamMembersName         string `json:"name"`
	// TeamMembersPhoneNumeber string `json:"phoneNumber"`
}

// type TeamType struct {
// 	TeamTypeDeprecated bool   `json:"deprecated"`
// 	TeamTypeHandle     string `json:"handle"`
// 	TeamTypeLabel      string `json:"label"`
// }

type TeamCollection struct {
	// Children    []string      `json:"children"`
	// Delete      bool          `json:"deleted"`
	// DisplayName string        `json:"displayName"`
	// Email       string        `json:"email"`
	// GroupName   string        `json:"groupName"`
	Id      int           `json:"id"`
	Members []TeamMembers `json:"members"`
	Name    string        `json:"name"`
	// Parent      string        `json:"parent"`
	// PhoneNumber string        `json:"phoneNumber"`
	// TeamTypeHandle TeamType      `json"teamType"`
}

func (c *Client) GetTeamSelection(ctx Context) ([]TeamCollection, error) {
	var v []TeamCollection

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams", nil)
	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, err
	}

	return v, err
}
