package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthUser struct {
	ID           int
	Firstname    string
	Surname      string
	Email        string
	Organisation string
	Roles        []string
	Locked       bool
	Suspended    bool
	Inactive     bool
}

type authUserResponse struct {
	ID        int      `json:"id"`
	Firstname string   `json:"firstname"`
	Surname   string   `json:"surname"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	Locked    bool     `json:"locked"`
	Suspended bool     `json:"suspended"`
	Inactive  bool     `json:"inactive"`
}

func (c *Client) User(ctx Context, id int) (AuthUser, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/auth/user/%d", id), nil)
	if err != nil {
		return AuthUser{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return AuthUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return AuthUser{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return AuthUser{}, newStatusError(resp)
	}

	var v authUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return AuthUser{}, err
	}

	user := AuthUser{
		ID:        v.ID,
		Firstname: v.Firstname,
		Surname:   v.Surname,
		Email:     v.Email,
		Locked:    v.Locked,
		Suspended: v.Suspended,
		Inactive:  v.Inactive,
	}

	for _, role := range v.Roles {
		if role == "OPG User" || role == "COP User" {
			user.Organisation = role
		} else {
			user.Roles = append(user.Roles, role)
		}
	}

	return user, err
}
