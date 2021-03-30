package sirius

import (
	"encoding/json"
	"net/http"
)

type ApiTaskTypes struct {
	Handle     string `json:"handle"`
	Incomplete string `json:"incomplete"`
	Category   string `json:"category"`
	Complete   string `json:"complete"`
	User       bool   `json:"user"`
}

type WholeTaskList struct {
	AllTaskList ApiTaskTypes `json:"task_types"`
}

func (c *Client) GetTaskDetails(ctx Context) (WholeTaskList, error) {
	var t WholeTaskList
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/tasktypes/supervision", nil)
	if err != nil {
		return t, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return t, err
	}
	defer resp.Body.Close()

	// io.Copy(os.Stdout, resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		return t, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return t, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return t, err
	}

	wholeTaskList := t

	return wholeTaskList, err
}
