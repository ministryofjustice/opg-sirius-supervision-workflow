package model

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestClient_GetStatus(t *testing.T) {
	status := func(s string) RefData { return RefData{Label: s} }

	tests := []struct {
		orders     []Order
		orderType  string
		wantStatus string
	}{
		{
			orders: []Order{
				{Status: status("Closed"), Type: "pfa"},
				{Status: status("Open"), Type: "pfa"},
				{Status: status("Duplicate"), Type: "pfa"},
				{Status: status("Active"), Type: "pfa"},
				{Status: status("Closed"), Type: "pfa"},
				{Status: status("Open"), Type: "pfa"},
				{Status: status("Duplicate"), Type: "pfa"},
			},
			wantStatus: "Active",
		},
		{
			orders: []Order{
				{Status: status("Closed"), Type: "pfa"},
				{Status: status("Open"), Type: "pfa"},
				{Status: status("Duplicate"), Type: "pfa"},
				{Status: status("Active"), Type: "pfa"},
				{Status: status("Closed"), Type: "hw"},
				{Status: status("Open"), Type: "hw"},
				{Status: status("Duplicate"), Type: "pfa"},
			},
			orderType:  "hw",
			wantStatus: "Open",
		},
		{
			orders: []Order{
				{Status: status("Open"), Type: "pfa"},
				{Status: status("Duplicate"), Type: "pfa"},
				{Status: status("Closed"), Type: "hw"},
				{Status: status("Open"), Type: "hw"},
				{Status: status("Duplicate"), Type: "pfa"},
			},
			wantStatus: "Open",
		},
		{
			orders: []Order{
				{Status: status("Duplicate"), Type: "pfa"},
				{Status: status("Closed"), Type: "hw"},
				{Status: status("Duplicate"), Type: "pfa"},
			},
			wantStatus: "Closed",
		},
		{
			orders: []Order{
				{Status: status("Duplicate"), Type: "pfa"},
			},
			wantStatus: "Duplicate",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			client := Client{Orders: test.orders}
			assert.Equal(t, test.wantStatus, client.GetStatus(test.orderType))
		})
	}
}

func TestClient_GetClosedCasesStatus(t *testing.T) {
	status := func(s string) RefData { return RefData{Label: s} }

	tests := []struct {
		orders     []Order
		orderType  string
		wantStatus string
	}{
		{
			orders: []Order{
				{Status: status("Closed"), Type: "pfa"},
				{Status: status("Open"), Type: "pfa"},
				{Status: status("Duplicate"), Type: "pfa"},
				{Status: status("Active"), Type: "pfa"},
				{Status: status("Closed"), Type: "pfa"},
				{Status: status("Open"), Type: "pfa"},
				{Status: status("Duplicate"), Type: "pfa"},
			},
			wantStatus: "Active",
		},
		{
			orders: []Order{
				{Status: status("Closed"), Type: "pfa"},
				{Status: status("Open"), Type: "pfa"},
				{Status: status("Duplicate"), Type: "pfa"},
				{Status: status("Active"), Type: "pfa"},
				{Status: status("Closed"), Type: "hw"},
				{Status: status("Open"), Type: "hw"},
				{Status: status("Duplicate"), Type: "pfa"},
			},
			orderType:  "hw",
			wantStatus: "Open",
		},
		{
			orders: []Order{
				{Status: status("Open"), Type: "pfa"},
				{Status: status("Duplicate"), Type: "pfa"},
				{Status: status("Closed"), Type: "hw"},
				{Status: status("Open"), Type: "hw"},
				{Status: status("Duplicate"), Type: "pfa"},
			},
			wantStatus: "Open",
		},
		{
			orders: []Order{
				{Status: status("Duplicate"), Type: "pfa"},
				{Status: status("Closed"), Type: "hw"},
				{Status: status("Duplicate"), Type: "pfa"},
			},
			wantStatus: "Duplicate",
		},
		{
			orders: []Order{
				{Status: status("Duplicate"), Type: "pfa"},
			},
			wantStatus: "Duplicate",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			client := Client{Orders: test.orders}
			assert.Equal(t, test.wantStatus, client.GetClosedCasesStatus(test.orderType))
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
	tests := []struct {
		name      string
		orders    []Order
		orderType string
		want      Order
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
					Status:         RefData{Handle: "CLOSED", Label: "Closed"},
				},
				{
					MadeActiveDate: NewDate("01/01/2019"),
					Status:         RefData{Handle: "ACTIVE", Label: "Active"},
				},
				{
					MadeActiveDate: NewDate("01/02/2020"),
					Date:           NewDate("01/02/2020"),
					Status:         RefData{Handle: "ACTIVE", Label: "Active"},
				},
			},
			want: Order{
				MadeActiveDate: NewDate("01/02/2020"),
				Date:           NewDate("01/02/2020"),
				Status:         RefData{Handle: "ACTIVE", Label: "Active"},
			},
		},
		{
			name: "It filters by order type if supplied",
			orders: []Order{
				{
					MadeActiveDate: NewDate("01/05/2020"),
					Date:           NewDate("01/08/2020"),
					Status:         RefData{Handle: "CLOSED", Label: "Closed"},
				},
				{
					MadeActiveDate: NewDate("01/01/2019"),
					Status:         RefData{Handle: "ACTIVE", Label: "Active"},
					Type:           "hw",
				},
				{
					MadeActiveDate: NewDate("01/02/2020"),
					Date:           NewDate("01/02/2020"),
					Status:         RefData{Handle: "ACTIVE", Label: "Active"},
				},
			},
			orderType: "hw",
			want: Order{
				MadeActiveDate: NewDate("01/01/2019"),
				Status:         RefData{Handle: "ACTIVE", Label: "Active"},
				Type:           "hw",
			},
		},
		{
			name: "It prioritises active orders with made active date over closed orders with a made active date",
			orders: []Order{
				{
					Date:   NewDate("01/08/2022"),
					Status: RefData{Handle: "ACTIVE", Label: "Active"},
				},
				{
					MadeActiveDate: NewDate("01/01/2023"),
					Status:         RefData{Handle: "CLOSED", Label: "Closed"},
				},
				{
					Date:           NewDate("01/06/2022"),
					MadeActiveDate: NewDate("01/05/2022"),
					Status:         RefData{Handle: "ACTIVE", Label: "Active"},
				},
			},
			want: Order{
				Date:           NewDate("01/06/2022"),
				MadeActiveDate: NewDate("01/05/2022"),
				Status:         RefData{Handle: "ACTIVE", Label: "Active"},
			},
		},
		{
			name:   "It returns no order if there are no orders",
			orders: []Order(nil),
			want:   Order{},
		},
	}
	for i, test := range tests {
		client := Client{}
		client.Orders = test.orders
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, client.GetMostRecentlyMadeActiveOrder(test.orderType))
		})
	}
}

func TestClient_GetActiveOrders(t *testing.T) {
	tests := []struct {
		orders    []Order
		orderType string
		want      []Order
	}{
		{
			orders: []Order{
				{
					Id:     1,
					Status: RefData{Handle: "CLOSED", Label: "Closed"},
				},
				{
					Id:     2,
					Status: RefData{Handle: "ACTIVE", Label: "Active"},
				},
				{
					Id:     3,
					Status: RefData{Handle: "DUPLICATE", Label: "Duplicate"},
				},
			},
			want: []Order{
				{
					Id:     2,
					Status: RefData{Handle: "ACTIVE", Label: "Active"},
				},
			},
		},
		{
			orders: []Order{
				{
					Id:     1,
					Status: RefData{Handle: "ACTIVE", Label: "Active"},
				},
				{
					Id:     2,
					Status: RefData{Handle: "ACTIVE", Label: "Active"},
				},
				{
					Id:     3,
					Status: RefData{Handle: "DUPLICATE", Label: "Duplicate"},
				},
			},
			want: []Order{
				{
					Id:     1,
					Status: RefData{Handle: "ACTIVE", Label: "Active"},
				},
				{
					Id:     2,
					Status: RefData{Handle: "ACTIVE", Label: "Active"},
				},
			},
		},
		{
			orders: []Order{
				{
					Id:     1,
					Status: RefData{Handle: "ACTIVE", Label: "Active"},
				},
				{
					Id:     2,
					Status: RefData{Handle: "ACTIVE", Label: "Active"},
					Type:   "hw",
				},
				{
					Id:     3,
					Status: RefData{Handle: "DUPLICATE", Label: "Duplicate"},
				},
			},
			orderType: "hw",
			want: []Order{
				{
					Id:     2,
					Status: RefData{Handle: "ACTIVE", Label: "Active"},
					Type:   "hw",
				},
			},
		},
		{
			orders: []Order{
				{
					Id:     1,
					Status: RefData{Handle: "CLOSED", Label: "Closed"},
				},
				{
					Id:     2,
					Status: RefData{Handle: "DUPLICATE", Label: "Duplicate"},
				},
			},
			want: []Order(nil),
		},
	}
	for i, test := range tests {
		client := Client{Orders: test.orders}
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, client.GetActiveOrders(test.orderType))
		})
	}
}

func TestClient_GetMostRecentOrder(t *testing.T) {
	tests := []struct {
		name      string
		orders    []Order
		orderType string
		want      Order
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
				MadeActiveDate: NewDate("01/05/2020"),
				Date:           NewDate("01/08/2020"),
			},
		},
		{
			name: "It filters by order type if supplied",
			orders: []Order{
				{
					MadeActiveDate: NewDate("01/05/2020"),
					Date:           NewDate("01/08/2020"),
				},
				{
					MadeActiveDate: NewDate("01/01/2019"),
					Type:           "pfa",
				},
				{
					MadeActiveDate: NewDate("01/02/2020"),
					Date:           NewDate("01/02/2020"),
				},
			},
			orderType: "pfa",
			want: Order{
				MadeActiveDate: NewDate("01/01/2019"),
				Type:           "pfa",
			},
		},
		{
			name: "It prioritises orders on date if no made active date",
			orders: []Order{
				{Date: NewDate("01/08/2020")},
				{Date: NewDate("01/01/2019")},
				{Date: NewDate("01/02/2020")},
			},
			want: Order{
				Date: NewDate("01/08/2020"),
			},
		},
		{
			name:   "It returns no order if there are no orders",
			orders: []Order{},
			want:   Order{},
		},
	}
	for i, test := range tests {
		client := Client{}
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, client.GetMostRecentOrder(test.orders, test.orderType))
		})
	}
}
