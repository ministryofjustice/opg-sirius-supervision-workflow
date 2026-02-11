package model

import (
    "fmt"
    "golang.org/x/text/language"
    "golang.org/x/text/message"
)

type Bond struct {
	Id                  int     `json:"id"`
	CourtRef            string  `json:"caseReferenceNumber"`
	FirstName           string  `json:"clientFirstName"`
	LastName            string  `json:"clientLastName"`
	CompanyName         string  `json:"companyName"`
	BondReferenceNumber string  `json:"bondReferenceNumber"`
	BondAmount          int     `json:"bondAmount"`
	BondIssuedDate      Date    `json:"bondIssuedDate"`
	BondClient          Client  `json:"client"`
	BondStatus          RefData `json:"bondStatus"`
}

func (b Bond) GetBondAmount() string {
    pounds := float64(b.BondAmount)/100
	p := message.NewPrinter(language.BritishEnglish)
	return p.Sprintf("£%.2f", pounds)
}
func (b Bond) GetURL() string {
	return fmt.Sprintf("/supervision/#/clients/%d", b.BondClient.Id)
}
