package sirius

import (
	"context"
	"net/http"
)

type UserDetails struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	PhoneNumber string          `json:"phoneNumber"`
	Teams       []MyDetailsTeam `json:"teams"`
	DisplayName string          `json:"displayName"`
	Deleted     bool            `json:"deleted"`
	Email       string          `json:"email"`
	Firstname   string          `json:"firstname"`
	Surname     string          `json:"surname"`
	Roles       []string        `json:"roles"`
	Locked      bool            `json:"locked"`
	Suspended   bool            `json:"suspended"`
}

type MyDetailsTeam struct {
	DisplayName string `json:"displayName"`
}

func (c *Client) SiriusUserDetails(ctx context.Context, cookies []*http.Cookie) (UserDetails, error) {
	myDetails := UserDetails{ID: 47,
		Name:        "system",
		PhoneNumber: "03004560300",
		Teams: []MyDetailsTeam{
			{DisplayName: "Allocations - (Supervision)"},
		},
		DisplayName: "system admin",
		Deleted:     false,
		Email:       "system.admin@opgtest.com",
		Firstname:   "system",
		Surname:     "admin",
		Roles:       []string{"OPG User", "System Admin"},
		Locked:      false,
		Suspended:   false}

	return myDetails, nil

	// req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url("/api/v1/users/current"), nil)
	// if err != nil {
	// 	return v, err
	// }
	// var xsrfToken string
	// for _, c := range cookies {
	// 	req.AddCookie(c)
	// 	if c.Name == "XSRF-TOKEN" {
	// 		xsrfToken = c.Value
	// 	}
	// }
	// req.Header.Add("OPG-Bypass-Membrane", "1")
	// req.Header.Add("X-XSRF-TOKEN", xsrfToken)

	// resp, err := c.http.Do(req)
	// if err != nil {
	// 	return v, err
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode == http.StatusUnauthorized {
	// 	return v, ErrUnauthorized
	// }

	// if resp.StatusCode != http.StatusOK {
	// 	return v, errors.New("returned non-2XX response")
	// }

	// err = json.NewDecoder(resp.Body).Decode(&v)
	// return v, err
}
