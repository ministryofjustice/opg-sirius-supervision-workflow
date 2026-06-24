package model

import (
	"fmt"
	"math"
	"strconv"
)

type DeputyImportantInformation struct {
	PanelDeputy bool `json:"panelDeputy"`
}

type DeputyAndCount struct {
	DeputyId int `json:"deputy"`
	Count    int `json:"count"`
}

type Deputy struct {
	Id                            int                        `json:"id"`
	DisplayName                   string                     `json:"displayName"`
	Type                          RefData                    `json:"deputyType"`
	Number                        int                        `json:"deputyNumber"`
	Address                       Address                    `json:"deputyAddress"`
	ExecutiveCaseManager          Assignee                   `json:"executiveCaseManager"`
	Assurance                     Assurance                  `json:"mostRecentlyCompletedAssurance"`
	ActiveClientCount             int                        `json:"activeClientCount"`
	ActiveNonCompliantClientCount int                        `json:"activeNonCompliantClientCount"`
	DeputyImportantInformation    DeputyImportantInformation `json:"deputyImportantInformation"`
	Firm                          Firm                       `json:"firm"`
}

func (d Deputy) IsSelected(selectedDeputies []string) bool {
	for _, a := range selectedDeputies {
		id, _ := strconv.Atoi(a)
		if d.Id == id {
			return true
		}
	}
	return false
}

func (d Deputy) GetURL() string {
	url := "/supervision/deputies/%d"
	if d.Type.Handle == "LAY" {
		url = "/supervision/#/deputy-hub/%d"
	}
	return fmt.Sprintf(url, d.Id)
}

func (d Deputy) GetFirm() Firm {
	return d.Firm
}

func (d Deputy) IsPro() bool {
	return d.Type.Handle == "PRO"
}

func (d Deputy) CalculateNonCompliance() string {
	if d.ActiveClientCount == 0 {
		return "0%"
	}
	percentage := (float64(d.ActiveNonCompliantClientCount) / float64(d.ActiveClientCount)) * 100
	return fmt.Sprintf("%.f%%", math.Round(percentage))
}

func (d Deputy) GetCountAsString(selectedDeputies []DeputyAndCount, urlPath string) string {
	for _, sd := range selectedDeputies {
		if d.Id == sd.DeputyId {
			stringValue := strconv.Itoa(sd.Count)
			return "(" + stringValue + ")"
		}
	}
	return "(0)"
}
