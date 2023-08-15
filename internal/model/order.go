package model

type Order struct {
	Id                     int          `json:"id"`
	Client                 Client       `json:"client"`
	Status                 RefData      `json:"orderStatus"`
	LatestAnnualReport     AnnualReport `json:"latestAnnualReport"`
	Date                   Date         `json:"orderDate"`
	MadeActiveDate         Date         `json:"madeActiveDate"`
	HowDeputyAppointed     RefData      `json:"howDeputyAppointed"`
	IntroductoryTargetDate Date         `json:"introductoryTargetDate"`
}

func (o Order) GetIntroductoryTargetDate() string {
	nullDate := NewDate("01/01/0001")
	if nullDate.Before(o.IntroductoryTargetDate) {
		return o.IntroductoryTargetDate.String()
	}
	return ""
}
