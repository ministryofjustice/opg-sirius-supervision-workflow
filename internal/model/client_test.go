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
