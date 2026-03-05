package model

import (
	"fmt"
	"sort"
	"strings"
)

type Bond struct {
	Id                  int               `json:"id"`
	CourtRef            string            `json:"caseReferenceNumber"`
	FirstName           string            `json:"clientFirstName"`
	LastName            string            `json:"clientLastName"`
	CompanyName         string            `json:"companyName"`
	BondReferenceNumber string            `json:"bondReferenceNumber"`
	BondAmount          int               `json:"bondAmount"` // amount in pounds
	BondIssuedDate      Date              `json:"bondIssuedDate"`
	BondClient          Client            `json:"client"`
	BondStatus          RefData           `json:"bondStatus"`
	Deputies          []string            `json:"deputyNames"`
}

func (b Bond) GetURL() string {
	return fmt.Sprintf("/supervision/#/clients/%d", b.BondClient.Id)
}

func (b Bond) GetDeputiesList() string {
	if len(b.Deputies) == 0 {
		return ""
	}

	names := make([]string, 0, len(b.Deputies))
	for _, name := range b.Deputies {

		if strings.TrimSpace(name) != "" {
			names = append(names, name)
		}
	}

	sort.Strings(names)
	return strings.Join(names, ", ")
}
