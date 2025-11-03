package model

import "fmt"

type Bond struct {
	Id                  int    `json:"id"`
	CourtRef            string `json:"caseReferenceNumber"`
	FirstName           string `json:"clientFirstName"`
	LastName            string `json:"clientLastName"`
	CompanyName         string `json:"companyName"`
	BondReferenceNumber string `json:"bondReferenceNumber"`
	BondAmount          int    `json:"bondAmount"`
	BondIssuedDate      Date   `json:"bondIssuedDate"`
}

func (b Bond) GetBondAmount() string {
	return "£" + fmt.Sprintf("%.2f", float64(b.BondAmount)/100)
}
