package model

type Order struct {
	Id                     int
	Client                 Client
	Type                   string
	Status                 RefData
	LatestAnnualReport     AnnualReport
	Date                   Date
	MadeActiveDate         Date
	HowDeputyAppointed     RefData
	IntroductoryTargetDate Date
}
