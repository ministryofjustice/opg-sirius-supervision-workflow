package sirius

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
)

type UserStatus string

func (us UserStatus) String() string {
	return string(us)
}

func (us UserStatus) TagColour() string {
	if us == "Suspended" {
		return "govuk-tag--grey"
	} else if us == "Locked" {
		return "govuk-tag--orange"
	} else {
		return ""
	}
}

type apiUser struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	Locked      bool   `json:"locked"`
	Suspended   bool   `json:"suspended"`
}

type apiUserList struct {
	Data []apiUser `json:"data"`
}

type User struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Status      UserStatus
}

func (c *Client) SearchUsers(ctx Context, search string) ([]User, error) {
	if len(search) < 3 {
		return nil, ClientError("Search term must be at least three characters")
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/api/search/users?query="+search, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v apiUserList
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	apiUsers := v.Data

	sort.SliceStable(apiUsers, func(i, j int) bool {
		if strings.EqualFold(apiUsers[i].Surname, apiUsers[j].Surname) {
			return strings.ToLower(apiUsers[i].DisplayName) < strings.ToLower(apiUsers[j].DisplayName)
		}

		return strings.ToLower(apiUsers[i].Surname) < strings.ToLower(apiUsers[j].Surname)
	})

	var users []User
	for _, u := range apiUsers {
		user := User{
			ID:          u.ID,
			DisplayName: u.DisplayName,
			Email:       u.Email,
			Status:      "Active",
		}

		if u.Suspended {
			user.Status = "Suspended"
		} else if u.Locked {
			user.Status = "Locked"
		}

		users = append(users, user)
	}

	return users, nil
}
