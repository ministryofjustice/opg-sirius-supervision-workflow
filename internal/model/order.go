package model

type Order struct {
	Id                 int          `json:"id"`
	Client             Client       `json:"client"`
	Status             RefData      `json:"orderStatus"`
	LatestAnnualReport AnnualReport `json:"latestAnnualReport"`
}
