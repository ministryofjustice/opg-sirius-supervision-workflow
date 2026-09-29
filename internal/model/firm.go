package model

import "fmt"

type Firm struct {
	Id     int
	Name   string
	Number int
}

func (f Firm) GetFirmURL() string {
	url := "/supervision/deputies/firm/%d"

	return fmt.Sprintf(url, f.Id)
}
