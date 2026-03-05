package model

import (
	"fmt"
)

type Bond struct {
	Id                  int     `json:"id"`
	CourtRef            string  `json:"caseReferenceNumber"`
	FirstName           string  `json:"clientFirstName"`
	LastName            string  `json:"clientLastName"`
	CompanyName         string  `json:"companyName"`
	BondReferenceNumber string  `json:"bondReferenceNumber"`
	BondAmount          int     `json:"bondAmount"` // amount in pounds
	BondIssuedDate      Date    `json:"bondIssuedDate"`
	BondClient          Client  `json:"client"`
	BondStatus          RefData `json:"bondStatus"`
}

func (b Bond) GetURL() string {
	return fmt.Sprintf("/supervision/#/clients/%d", b.BondClient.Id)
}
