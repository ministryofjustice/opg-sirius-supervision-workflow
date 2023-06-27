package model

type Client struct {
	Id                   int      `json:"id"`
	CaseRecNumber        string   `json:"caseRecNumber"`
	FirstName            string   `json:"firstname"`
	Surname              string   `json:"surname"`
	SupervisionCaseOwner Assignee `json:"supervisionCaseOwner"`
	FeePayer             Deputy   `json:"feePayer"`
	Orders               []Order  `json:"cases"`
	SupervisionLevel     string   `json:"supervisionLevel"`
}

func (c Client) GetReportDueDate() string {
	if len(c.Orders) > 0 {
		return c.Orders[0].LatestAnnualReport.DueDate
	}
	return ""
}

func (c Client) GetStatus() string {
	orderStatuses := make(map[string]string)

	for _, order := range c.Orders {
		label := order.Status.Label
		orderStatuses[label] = label
	}

	statuses := []string{"Active", "Open", "Closed", "Duplicate"}
	for _, status := range statuses {
		if _, found := orderStatuses[status]; found {
			return status
		}
	}

	return ""
}
