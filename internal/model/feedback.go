package model

type FeedbackForm struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	CaseNumber string `json:"caseNumber"`
	Message    string `json:"message"`
}
