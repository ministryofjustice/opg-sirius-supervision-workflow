package model

type Assurance struct {
	ReportReviewDate Date    `json:"reportReviewDate"`
	ReportMarkedAs   RefData `json:"reportMarkedAs"`
	Type             RefData `json:"assuranceType"`
}

type RAGRating struct {
	Name   string
	Colour string
}

func (a Assurance) IsPDR() bool {
	return a.Type.Handle == "PDR"
}

func (a Assurance) GetRAGRating() RAGRating {
	var rag RAGRating
	switch a.ReportMarkedAs.Handle {
	case "RED":
		rag.Name = "High risk"
		rag.Colour = "red"
	case "AMBER":
		rag.Name = "Medium risk"
		rag.Colour = "orange"
	case "GREEN":
		rag.Name = "Low risk"
		rag.Colour = "green"
	}
	return rag
}
