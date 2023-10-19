package model

import "fmt"

type Deputy struct {
	Id                            int       `json:"id"`
	DisplayName                   string    `json:"displayName"`
	Type                          RefData   `json:"deputyType"`
	Number                        int       `json:"deputyNumber"`
	Address                       Address   `json:"deputyAddress"`
	ExecutiveCaseManager          Assignee  `json:"executiveCaseManager"`
	Assurance                     Assurance `json:"mostRecentlyCompletedAssurance"`
	ActiveClientCount             int       `json:"activeClientCount"`
	ActiveNonCompliantClientCount int       `json:"activeNonCompliantClientCount"`
}

func (d Deputy) GetURL() string {
	url := "/supervision/deputies/%d"
	if d.Type.Handle == "LAY" {
		url = "/supervision/#/deputy-hub/%d"
	}
	return fmt.Sprintf(url, d.Id)
}

func (d Deputy) IsPro() bool {
	return d.Type.Handle == "PRO"
}

func (d Deputy) CalculateNonCompliance() string {
	if d.ActiveClientCount == 0 {
		return "0%"
	}
	return fmt.Sprintf("%.f%%", (float64(d.ActiveNonCompliantClientCount)/float64(d.ActiveClientCount))*100)
}
