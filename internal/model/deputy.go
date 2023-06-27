package model

import "fmt"

type Deputy struct {
	Id          int     `json:"id"`
	DisplayName string  `json:"displayName"`
	Type        RefData `json:"deputyType"`
}

func (d Deputy) GetURL() string {
	url := "/supervision/deputies/%d"
	if d.Type.Handle == "LAY" {
		url = "/supervision/#/deputy-hub/%d"
	}
	return fmt.Sprintf(url, d.Id)
}
