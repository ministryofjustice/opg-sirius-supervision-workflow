package model

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestDeputy_GetURL(t *testing.T) {
	tests := []struct {
		name        string
		deputyType  string
		expectedUrl string
	}{
		{
			name:        "Professional deputy URL",
			deputyType:  "PRO",
			expectedUrl: "/supervision/deputies/13",
		},
		{
			name:        "PA deputy URL",
			deputyType:  "PA",
			expectedUrl: "/supervision/deputies/13",
		},
		{
			name:        "Lay deputy URL",
			deputyType:  "LAY",
			expectedUrl: "/supervision/#/deputy-hub/13",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			deputy := Deputy{Id: 13, Type: RefData{Handle: test.deputyType}}
			assert.Equal(t, test.expectedUrl, deputy.GetURL())
		})
	}
}

func TestDeputy_IsPro(t *testing.T) {
	tests := []struct {
		deputy    Deputy
		wantIsPro bool
	}{
		{
			deputy:    Deputy{},
			wantIsPro: false,
		},
		{
			deputy:    Deputy{Type: RefData{Handle: "PRO"}},
			wantIsPro: true,
		},
		{
			deputy:    Deputy{Type: RefData{Handle: "PA"}},
			wantIsPro: false,
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.wantIsPro, test.deputy.IsPro())
		})
	}
}

func TestDeputy_CalculateNonCompliance(t *testing.T) {
	tests := []struct {
		activeClientCount int
		nonCompliantCount int
		want              string
	}{
		{
			activeClientCount: 0,
			nonCompliantCount: 0,
			want:              "0%",
		},
		{
			activeClientCount: 10,
			nonCompliantCount: 0,
			want:              "0%",
		},
		{
			activeClientCount: 0,
			nonCompliantCount: 10,
			want:              "0%",
		},
		{
			activeClientCount: 1,
			nonCompliantCount: 1,
			want:              "100%",
		},
		{
			activeClientCount: 3,
			nonCompliantCount: 1,
			want:              "33%",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			deputy := Deputy{
				ActiveClientCount:             test.activeClientCount,
				ActiveNonCompliantClientCount: test.nonCompliantCount,
			}
			assert.Equal(t, test.want, deputy.CalculateNonCompliance())
		})
	}
}
