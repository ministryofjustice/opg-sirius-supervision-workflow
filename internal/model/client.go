package model

import "fmt"

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

func (c Client) GetURL() string {
	return fmt.Sprintf("/supervision/#/clients/%d", c.Id)
}

func (c Client) GetActiveOrders() []Order {
	var activeOrders []Order
	for _, order := range c.Orders {
		if order.Status.Handle == "ACTIVE" {
			activeOrders = append(activeOrders, order)
		}
	}
	return activeOrders
}

func (c Client) GetMostRecentOrder(orders []Order) Order {
	var mostRecentlyMadeActiveOrder Order
	var mostRecentlyMadeOrder Order

	for _, order := range orders {
		if order.MadeActiveDate.After(mostRecentlyMadeActiveOrder.MadeActiveDate) {
			mostRecentlyMadeActiveOrder = order
		}
		if order.Date.After(mostRecentlyMadeOrder.Date) {
			mostRecentlyMadeOrder = order
		}
	}
	if mostRecentlyMadeActiveOrder.MadeActiveDate.IsNull() {
		return mostRecentlyMadeOrder
	}
	return mostRecentlyMadeActiveOrder
}

func (c Client) GetMostRecentlyMadeActiveOrder() Order {
	activeOrders := c.GetActiveOrders()
	if len(activeOrders) > 0 {
		return c.GetMostRecentOrder(activeOrders)
	}
	return c.GetMostRecentOrder(c.Orders)
}
