package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) EditTeam(ctx Context, team Team) error {
	form := url.Values{
		"name":             {team.DisplayName},
		"email":            {team.Email},
		"phoneNumber":      {team.PhoneNumber},
		"teamType[handle]": {team.Type},
	}

	for i, member := range team.Members {
		form.Add(fmt.Sprintf("members[%d][id]", i), strconv.Itoa(member.ID))
	}

	body := strings.NewReader(form.Encode())

	requestURL := fmt.Sprintf("/api/team/%d", team.ID)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		var v struct {
			Data struct {
				Errors ValidationErrors `json:"errorMessages"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return &ValidationError{
				Errors: v.Data.Errors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}
