package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) EditMyDetails(ctx Context, id int, phoneNumber string) error {
	var v struct {
		Detail           string           `json:"detail"`
		ValidationErrors ValidationErrors `json:"validation_errors"`
	}

	var body = strings.NewReader("{\"phoneNumber\":\"" + phoneNumber + "\"}")

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/users/%d/updateTelephoneNumber", id), body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return &ValidationError{
				Message: v.Detail,
				Errors:  v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}
