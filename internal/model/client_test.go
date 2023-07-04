package model

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestClient_GetStatus(t *testing.T) {
	tests := []struct {
		orderStatuses []string
		wantStatus    string
	}{
		{
			orderStatuses: []string{"Closed", "Open", "Duplicate", "Active", "Closed", "Open", "Duplicate"},
			wantStatus:    "Active",
		},
		{
			orderStatuses: []string{"Open", "Duplicate", "Closed", "Open", "Duplicate"},
			wantStatus:    "Open",
		},
		{
			orderStatuses: []string{"Duplicate", "Closed", "Duplicate"},
			wantStatus:    "Closed",
		},
		{
			orderStatuses: []string{"Duplicate"},
			wantStatus:    "Duplicate",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			var client Client
			for _, status := range test.orderStatuses {
				client.Orders = append(client.Orders, Order{Status: RefData{Label: status}})
			}
			assert.Equal(t, test.wantStatus, client.GetStatus())
		})
	}
}

func TestClient_GetReportDueDate(t *testing.T) {
	client := Client{
		Orders: []Order{
			{
				LatestAnnualReport: AnnualReport{
					DueDate: "12/02/2020",
				},
			},
		},
	}
	assert.Equal(t, "12/02/2020", client.GetReportDueDate())
}

func TestClient_GetURL(t *testing.T) {
	assert.Equal(t, "/supervision/#/clients/0", Client{}.GetURL())
	assert.Equal(t, "/supervision/#/clients/12", Client{Id: 12}.GetURL())
}

func TestClient_GetMostRecentlyMadeActiveOrder(t *testing.T) {
	client := Client{
		Orders: []Order{
			{MadeActiveDate: NewDate("01/12/2022")},
			{MadeActiveDate: NewDate("01/06/2023")},
			{MadeActiveDate: NewDate("01/07/2021")},
		},
	}
	assert.Equal(t, "01/06/2023", client.GetMostRecentlyMadeActiveOrder().MadeActiveDate.String())
	assert.Equal(t, Order{}, Client{}.GetMostRecentlyMadeActiveOrder())
}
