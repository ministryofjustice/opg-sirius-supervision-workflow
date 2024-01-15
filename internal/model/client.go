package model

import (
	"fmt"
)

type Client struct {
	Id                   int       `json:"id"`
	CaseRecNumber        string    `json:"caseRecNumber"`
	FirstName            string    `json:"firstname"`
	Surname              string    `json:"surname"`
	SupervisionCaseOwner Assignee  `json:"supervisionCaseOwner"`
	FeePayer             Deputy    `json:"feePayer"`
	Orders               []Order   `json:"cases"`
	SupervisionLevel     RefData   `json:"supervisionLevel"`
	ActiveCaseType       RefData   `json:"activeCaseType"`
	HWDeputyTypes        []RefData `json:"hwDeputyTypes"`
	LastActionDate       Date      `json:"lastActionDate"`
	CachedDebtTotal      float64   `json:"cachedDebtTotal"`
	ClosedOnDate         Date      `json:"closedOnDate"`
}

func (c Client) GetReportDueDate() string {
	if len(c.Orders) > 0 {
		return c.Orders[0].LatestAnnualReport.DueDate
	}
	return ""
}

func (c Client) GetStatus(orderType string, closedCases bool) string {
	orderStatuses := make(map[string]string)
	var statuses []string

	for _, order := range c.Orders {
		if orderType == "" || orderType == order.Type {
			label := order.Status.Label
			orderStatuses[label] = label
		}
	}
	if closedCases {
		statuses = []string{"Active", "Open", "Duplicate", "Closed"}
	} else {
		statuses = []string{"Active", "Open", "Closed", "Duplicate"}
	}
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

func (c Client) GetCachedDebtTotal() string {
	if c.CachedDebtTotal <= 0 {
		return "-"
	}
	return "£" + fmt.Sprintf("%.2f", c.CachedDebtTotal/100)
}

func (c Client) GetActiveOrders(orderType string) []Order {
	var activeOrders []Order
	for _, order := range c.Orders {
		if order.Status.Handle == "ACTIVE" && (orderType == "" || order.Type == orderType) {
			activeOrders = append(activeOrders, order)
		}
	}
	return activeOrders
}

func (c Client) GetMostRecentOrder(orders []Order, orderType string) Order {
	var mostRecentlyMadeActiveOrder Order
	var mostRecentlyMadeOrder Order

	for _, order := range orders {
		if orderType != "" && order.Type != orderType {
			continue
		}
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

func (c Client) GetMostRecentlyMadeActiveOrder(orderType string) Order {
	activeOrders := c.GetActiveOrders(orderType)
	if len(activeOrders) > 0 {
		return c.GetMostRecentOrder(activeOrders, orderType)
	}
	return c.GetMostRecentOrder(c.Orders, orderType)
}
