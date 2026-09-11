package model

import (
	"fmt"
	"sort"
	"strings"
)

type Bond struct {
	Id                  int
	CourtRef            string
	FirstName           string
	LastName            string
	CompanyName         string
	BondReferenceNumber string
	BondAmount          int // amount in pounds
	BondIssuedDate      Date
	BondClient          Client
	BondStatus          RefData
	Deputies            []string
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
