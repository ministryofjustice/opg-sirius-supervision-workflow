package model

import (
	"fmt"
	"math"
	"strconv"
)

type DeputyImportantInformation struct {
	PanelDeputy bool
}

type Deputy struct {
	Id                            int
	DisplayName                   string
	Type                          RefData
	Number                        int
	Address                       Address
	ExecutiveCaseManager          Assignee
	Assurance                     Assurance
	ActiveClientCount             int
	ActiveNonCompliantClientCount int
	DeputyImportantInformation    DeputyImportantInformation
	Firm                          Firm
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
