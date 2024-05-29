package model

import "fmt"

type Firm struct {
	Id     int    `json:"id"`
	Name   string `json:"firmName"`
	Number int    `json:"firmNumber"`
}

func (f Firm) GetFirmURL() string {
	url := "/supervision/deputies/firm/%d"

	return fmt.Sprintf(url, f.Id)
}
