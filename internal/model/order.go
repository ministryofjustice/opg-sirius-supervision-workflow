package model

type Order struct {
	Id                     int          `json:"id"`
	Client                 Client       `json:"client"`
	Type                   string       `json:"caseSubtype"`
	Status                 RefData      `json:"orderStatus"`
	LatestAnnualReport     AnnualReport `json:"latestAnnualReport"`
	Date                   Date         `json:"orderDate"`
	MadeActiveDate         Date         `json:"madeActiveDate"`
	HowDeputyAppointed     RefData      `json:"howDeputyAppointed"`
	IntroductoryTargetDate Date         `json:"introductoryTargetDate"`
}
