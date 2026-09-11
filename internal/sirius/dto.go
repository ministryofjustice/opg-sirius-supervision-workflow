package sirius

import (
	"strconv"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
)

type pageInformationResponse struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

func (r pageInformationResponse) model() model.PageInformation {
	return model.PageInformation{
		PageCurrent: r.Current,
		PageTotal:   r.Total,
	}
}

type assigneeAndCountResponse struct {
	AssigneeId int `json:"assignee"`
	Count      int `json:"count"`
}

func (r assigneeAndCountResponse) model() model.AssigneeAndCount {
	return model.AssigneeAndCount{
		AssigneeId: r.AssigneeId,
		Count:      r.Count,
	}
}

func assigneeAndCountsModel(rs []assigneeAndCountResponse) []model.AssigneeAndCount {
	if rs == nil {
		return nil
	}

	assigneeAndCounts := make([]model.AssigneeAndCount, 0, len(rs))
	for _, r := range rs {
		assigneeAndCounts = append(assigneeAndCounts, r.model())
	}
	return assigneeAndCounts
}

type refDataResponse struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (r *refDataResponse) model() model.RefData {
	if r == nil {
		return model.RefData{}
	}

	return model.RefData{
		Handle: r.Handle,
		Label:  r.Label,
	}
}

type teamResponse struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

func (r teamResponse) model() model.Team {
	return model.Team{
		Id:   r.ID,
		Name: r.DisplayName,
	}
}

type assigneeResponse struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

func (r *assigneeResponse) model() model.Assignee {
	if r == nil {
		return model.Assignee{}
	}

	return model.Assignee{
		Id:   r.ID,
		Name: r.DisplayName,
	}
}

type assigneeWithTeamsResponse struct {
	ID          int            `json:"id"`
	Teams       []teamResponse `json:"teams,omitempty"`
	DisplayName string         `json:"displayName"`
}

func (r *assigneeWithTeamsResponse) model() model.Assignee {
	if r == nil {
		return model.Assignee{}
	}

	var teams []model.Team
	if r.Teams != nil {
		teams = make([]model.Team, 0, len(r.Teams))
		for _, team := range r.Teams {
			teams = append(teams, team.model())
		}
	}

	return model.Assignee{
		Id:    r.ID,
		Name:  r.DisplayName,
		Teams: teams,
	}
}

type annualReportResponse struct {
	DueDate string `json:"dueDate"`
}

func (r *annualReportResponse) model() model.AnnualReport {
	if r == nil {
		return model.AnnualReport{}
	}

	return model.AnnualReport{
		DueDate: r.DueDate,
	}
}

type addressResponse struct {
	Town string `json:"town"`
}

func (r *addressResponse) model() model.Address {
	if r == nil {
		return model.Address{}
	}

	return model.Address{
		Town: r.Town,
	}
}

type assuranceResponse struct {
	ReportReviewDate string           `json:"reportReviewDate" pact:"example=2023-01-01T00:00:00+00:00"`
	ReportMarkedAs   *refDataResponse `json:"reportMarkedAs,omitempty"`
	AssuranceType    *refDataResponse `json:"assuranceType,omitempty"`
}

func (r *assuranceResponse) model() (model.Assurance, error) {
	if r == nil {
		return model.Assurance{}, nil
	}

	reportReviewDate, err := parseModelDate(r.ReportReviewDate)
	if err != nil {
		return model.Assurance{}, err
	}

	return model.Assurance{
		ReportReviewDate: reportReviewDate,
		ReportMarkedAs:   r.ReportMarkedAs.model(),
		Type:             r.AssuranceType.model(),
	}, nil
}

type deputyImportantInformationResponse struct {
	PanelDeputy bool `json:"panelDeputy"`
}

func (r *deputyImportantInformationResponse) model() model.DeputyImportantInformation {
	if r == nil {
		return model.DeputyImportantInformation{}
	}

	return model.DeputyImportantInformation{
		PanelDeputy: r.PanelDeputy,
	}
}

type firmResponse struct {
	ID         int    `json:"id"`
	FirmName   string `json:"firmName"`
	FirmNumber int    `json:"firmNumber"`
}

func (r *firmResponse) model() model.Firm {
	if r == nil {
		return model.Firm{}
	}

	return model.Firm{
		Id:     r.ID,
		Name:   r.FirmName,
		Number: r.FirmNumber,
	}
}

type deputySummaryResponse struct {
	ID          int              `json:"id"`
	DisplayName string           `json:"displayName"`
	DeputyType  *refDataResponse `json:"deputyType,omitempty"`
	Firm        *firmResponse    `json:"firm,omitempty"`
}

func (r *deputySummaryResponse) model() model.Deputy {
	if r == nil {
		return model.Deputy{}
	}

	return model.Deputy{
		Id:          r.ID,
		DisplayName: r.DisplayName,
		Type:        r.DeputyType.model(),
		Firm:        r.Firm.model(),
	}
}

type deputyResponse struct {
	ID                             int                                 `json:"id"`
	DeputyNumber                   int                                 `json:"deputyNumber"`
	DisplayName                    string                              `json:"displayName"`
	DeputyType                     *refDataResponse                    `json:"deputyType,omitempty"`
	DeputyAddress                  *addressResponse                    `json:"deputyAddress,omitempty"`
	ExecutiveCaseManager           *assigneeResponse                   `json:"executiveCaseManager,omitempty"`
	MostRecentlyCompletedAssurance *assuranceResponse                  `json:"mostRecentlyCompletedAssurance,omitempty"`
	ActiveClientCount              int                                 `json:"activeClientCount"`
	ActiveNonCompliantClientCount  int                                 `json:"activeNonCompliantClientCount"`
	DeputyImportantInformation     *deputyImportantInformationResponse `json:"deputyImportantInformation,omitempty"`
	Firm                           *firmResponse                       `json:"firm,omitempty"`
}

func (r deputyResponse) model() (model.Deputy, error) {
	assurance, err := r.MostRecentlyCompletedAssurance.model()
	if err != nil {
		return model.Deputy{}, err
	}

	return model.Deputy{
		Id:                            r.ID,
		DisplayName:                   r.DisplayName,
		Type:                          r.DeputyType.model(),
		Number:                        r.DeputyNumber,
		Address:                       r.DeputyAddress.model(),
		ExecutiveCaseManager:          r.ExecutiveCaseManager.model(),
		Assurance:                     assurance,
		ActiveClientCount:             r.ActiveClientCount,
		ActiveNonCompliantClientCount: r.ActiveNonCompliantClientCount,
		DeputyImportantInformation:    r.DeputyImportantInformation.model(),
		Firm:                          r.Firm.model(),
	}, nil
}

type clientSummaryResponse struct {
	ID                   int                        `json:"id"`
	UID                  string                     `json:"uId,omitempty"`
	CaseRecNumber        string                     `json:"caseRecNumber"`
	Salutation           string                     `json:"salutation,omitempty"`
	FirstName            string                     `json:"firstname"`
	MiddleNames          string                     `json:"middlenames,omitempty"`
	Surname              string                     `json:"surname"`
	SupervisionCaseOwner *assigneeWithTeamsResponse `json:"supervisionCaseOwner,omitempty"`
	FeePayer             *deputySummaryResponse     `json:"feePayer,omitempty"`
}

func (r *clientSummaryResponse) model() model.Client {
	if r == nil {
		return model.Client{}
	}

	return model.Client{
		Id:                   r.ID,
		CaseRecNumber:        r.CaseRecNumber,
		FirstName:            r.FirstName,
		Surname:              r.Surname,
		SupervisionCaseOwner: r.SupervisionCaseOwner.model(),
		FeePayer:             r.FeePayer.model(),
	}
}

type orderResponse struct {
	ID                     int                    `json:"id"`
	UID                    string                 `json:"uId,omitempty"`
	Client                 *clientSummaryResponse `json:"client,omitempty"`
	CaseRecNumber          string                 `json:"caseRecNumber,omitempty"`
	CaseType               string                 `json:"caseType,omitempty"`
	Type                   string                 `json:"caseSubtype"`
	OrderStatus            *refDataResponse       `json:"orderStatus,omitempty"`
	LatestAnnualReport     *annualReportResponse  `json:"latestAnnualReport,omitempty"`
	Date                   string                 `json:"orderDate" pact:"example=2023-01-01T00:00:00+00:00"`
	MadeActiveDate         string                 `json:"madeActiveDate" pact:"example=2023-01-01T00:00:00+00:00"`
	HowDeputyAppointed     *refDataResponse       `json:"howDeputyAppointed,omitempty"`
	IntroductoryTargetDate string                 `json:"introductoryTargetDate" pact:"example=2023-01-01T00:00:00+00:00"`
}

func (r orderResponse) model() (model.Order, error) {
	orderDate, err := parseModelDate(r.Date)
	if err != nil {
		return model.Order{}, err
	}

	madeActiveDate, err := parseModelDate(r.MadeActiveDate)
	if err != nil {
		return model.Order{}, err
	}

	introductoryTargetDate, err := parseModelDate(r.IntroductoryTargetDate)
	if err != nil {
		return model.Order{}, err
	}

	return model.Order{
		Id:                     r.ID,
		Client:                 r.Client.model(),
		Type:                   r.Type,
		Status:                 r.OrderStatus.model(),
		LatestAnnualReport:     r.LatestAnnualReport.model(),
		Date:                   orderDate,
		MadeActiveDate:         madeActiveDate,
		HowDeputyAppointed:     r.HowDeputyAppointed.model(),
		IntroductoryTargetDate: introductoryTargetDate,
	}, nil
}

type clientResponse struct {
	ID                   int                        `json:"id"`
	CaseRecNumber        string                     `json:"caseRecNumber"`
	FirstName            string                     `json:"firstname"`
	Surname              string                     `json:"surname"`
	SupervisionCaseOwner *assigneeWithTeamsResponse `json:"supervisionCaseOwner,omitempty"`
	FeePayer             *deputySummaryResponse     `json:"feePayer,omitempty"`
	Orders               []orderResponse            `json:"cases,omitempty"`
	SupervisionLevel     *refDataResponse           `json:"supervisionLevel,omitempty"`
	ActiveCaseType       *refDataResponse           `json:"activeCaseType,omitempty"`
	DeputyTypes          []refDataResponse          `json:"deputyTypes,omitempty"`
	LastActionDate       string                     `json:"lastActionDate" pact:"example=2023-01-01T00:00:00+00:00"`
	CachedDebtTotal      float64                    `json:"cachedDebtTotal"`
	ClosedOnDate         string                     `json:"closedOnDate" pact:"example=2023-01-01T00:00:00+00:00"`
}

func (r clientResponse) model() (model.Client, error) {
	var orders []model.Order
	if r.Orders != nil {
		orders = make([]model.Order, 0, len(r.Orders))
		for _, order := range r.Orders {
			mappedOrder, err := order.model()
			if err != nil {
				return model.Client{}, err
			}
			orders = append(orders, mappedOrder)
		}
	}

	lastActionDate, err := parseModelDate(r.LastActionDate)
	if err != nil {
		return model.Client{}, err
	}

	closedOnDate, err := parseModelDate(r.ClosedOnDate)
	if err != nil {
		return model.Client{}, err
	}

	var deputyTypes []model.RefData
	if r.DeputyTypes != nil {
		deputyTypes = make([]model.RefData, 0, len(r.DeputyTypes))
		for _, deputyType := range r.DeputyTypes {
			deputyTypes = append(deputyTypes, deputyType.model())
		}
	}

	return model.Client{
		Id:                   r.ID,
		CaseRecNumber:        r.CaseRecNumber,
		FirstName:            r.FirstName,
		Surname:              r.Surname,
		SupervisionCaseOwner: r.SupervisionCaseOwner.model(),
		FeePayer:             r.FeePayer.model(),
		Orders:               orders,
		SupervisionLevel:     r.SupervisionLevel.model(),
		ActiveCaseType:       r.ActiveCaseType.model(),
		DeputyTypes:          deputyTypes,
		LastActionDate:       lastActionDate,
		CachedDebtTotal:      r.CachedDebtTotal,
		ClosedOnDate:         closedOnDate,
	}, nil
}

type reassignResponse struct {
	ReassignName string `json:"reassignName"`
}

func parseModelDate(value string) (model.Date, error) {
	var date model.Date
	if value == "" {
		return date, nil
	}

	err := date.UnmarshalJSON([]byte(strconv.Quote(value)))
	return date, err
}
