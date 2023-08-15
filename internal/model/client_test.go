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

func TestClient_GetMostRecentlyMadeActiveOrderPrioritisesOrdersCorrectly(t *testing.T) {
	tests := []struct {
		name   string
		orders []Order
		want   Order
	}{
		{
			name: "It prioritises orders with made active date over date",
			orders: []Order{
				{
					MadeActiveDate: NewDate("01/05/2020"),
					Date:           NewDate("01/08/2020"),
				},
				{
					MadeActiveDate: NewDate("01/01/2019"),
				},
				{
					MadeActiveDate: NewDate("01/02/2020"),
					Date:           NewDate("01/02/2020"),
				},
			},
			want: Order{
				Date:           NewDate("01/08/2020"),
				MadeActiveDate: NewDate("01/05/2020"),
			},
		},
		{
			name: "It prioritises active orders",
			orders: []Order{
				{
					MadeActiveDate: NewDate("01/05/2020"),
					Date:           NewDate("01/08/2020"),
					Status: RefData{
						Handle: "CLOSED",
						Label:  "Closed",
					},
				},
				{
					MadeActiveDate: NewDate("01/01/2019"),
					Status: RefData{
						Handle: "ACTIVE",
						Label:  "Active",
					},
				},
				{
					MadeActiveDate: NewDate("01/02/2020"),
					Date:           NewDate("01/02/2020"),
					Status: RefData{
						Handle: "ACTIVE",
						Label:  "Active",
					},
				},
			},
			want: Order{
				MadeActiveDate: NewDate("01/02/2020"),
				Date:           NewDate("01/02/2020"),
				Status: RefData{
					Handle: "ACTIVE",
					Label:  "Active",
				},
			},
		},
		{
			name: "It prioritises any order with made active date over other orders with a date",
			orders: []Order{
				{
					MadeActiveDate: NewDate("01/05/2020"),
					Date:           NewDate("01/08/2023"),
				},
				{
					MadeActiveDate: NewDate("01/01/2023"),
				},
				{
					MadeActiveDate: NewDate("01/09/2022"),
					Date:           NewDate("01/10/2022"),
				},
			},
			want: Order{
				MadeActiveDate: NewDate("01/01/2023"),
			},
		},
		{
			name: "It prioritises active orders with made active date over closed orders with a made active date",
			orders: []Order{
				{
					Date: NewDate("01/08/2022"),
					Status: RefData{
						Handle: "ACTIVE",
						Label:  "Active",
					},
				},
				{
					MadeActiveDate: NewDate("01/01/2023"),
					Status: RefData{
						Handle: "CLOSED",
						Label:  "Closed",
					},
				},
				{
					Date:           NewDate("01/06/2022"),
					MadeActiveDate: NewDate("01/05/2022"),
					Status: RefData{
						Handle: "ACTIVE",
						Label:  "Active",
					},
				},
			},
			want: Order{
				Date:           NewDate("01/06/2022"),
				MadeActiveDate: NewDate("01/05/2022"),
				Status: RefData{
					Handle: "ACTIVE",
					Label:  "Active",
				},
			},
		},
	}
	for i, test := range tests {
		client := Client{}
		client.Orders = test.orders
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, client.GetMostRecentlyMadeActiveOrder())
		})
	}
}

func TestClient_GetActiveOrders(t *testing.T) {
	tests := []struct {
		orders []Order
		want   []Order
	}{
		{
			orders: []Order{
				{
					Id: 1,
					Status: RefData{
						Handle: "CLOSED",
						Label:  "Closed",
					},
				},
				{
					Id: 2,
					Status: RefData{
						Handle: "ACTIVE",
						Label:  "Active",
					},
				},
				{
					Id: 3,
					Status: RefData{
						Handle: "DUPLICATE",
						Label:  "Duplicate",
					},
				},
			},
			want: []Order{
				{
					Id: 1,
					Status: RefData{
						Handle: "CLOSED",
						Label:  "Closed",
					},
				},
			},
		},
		{
			orders: []Order{
				{
					Id: 1,
					Status: RefData{
						Handle: "ACTIVE",
						Label:  "Active",
					},
				},
				{
					Id: 2,
					Status: RefData{
						Handle: "ACTIVE",
						Label:  "Active",
					},
				},
				{
					Id: 3,
					Status: RefData{
						Handle: "DUPLICATE",
						Label:  "Duplicate",
					},
				},
			},
			want: []Order{
				{
					Id: 1,
					Status: RefData{
						Handle: "ACTIVE",
						Label:  "Active",
					},
				},
				{
					Id: 2,
					Status: RefData{
						Handle: "ACTIVE",
						Label:  "Active",
					},
				},
			},
		},
	}
	for i, test := range tests {
		client := Client{}
		client.Orders = test.orders
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, client.GetMostRecentlyMadeActiveOrder())
		})
	}
}
