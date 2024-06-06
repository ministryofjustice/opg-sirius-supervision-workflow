package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"net/http"
)

func (c *ApiClient) SubmitFeedback(ctx Context, form model.FeedbackForm) error {
	var body bytes.Buffer
	var err error
	fmt.Print("Vars")
	fmt.Println(form)

	err = json.NewEncoder(&body).Encode(form)

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, "/api/supervision-feedback", &body)

	if err != nil {
		c.logErrorRequest(req, err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		c.logResponse(req, resp, err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			c.logResponse(req, resp, err)
			return &ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}
