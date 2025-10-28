package model

import "fmt"

type Bond struct {
	Id                  int     `json:"id"`
	CourtRef            string  `json:"caseReferenceNumber"`
	FirstName           string  `json:"clientFirstName"`
	LastName            string  `json:"clientLastName"`
	CompanyName         string  `json:"companyName"`
	BondReferenceNumber string  `json:"bondReferenceNumber"`
	BondAmount          float32 `json:"bondAmount"`
	BondIssuedDate      Date    `json:"bondIssuedDate"`
}

func (b Bond) GetBondAmount() string {
	return "£" + fmt.Sprintf("%.2f", b.BondAmount)
}
