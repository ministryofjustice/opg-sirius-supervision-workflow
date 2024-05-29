package model

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestFirm_GetURL(t *testing.T) {
	tests := []struct {
		firm        Firm
		expectedUrl string
	}{
		{
			firm: Firm{
				Id:     65345,
				Name:   "firm test 1",
				Number: 9999,
			},
			expectedUrl: "/supervision/deputies/firm/65345",
		},
		{
			firm: Firm{
				Id:     9952223,
				Name:   "firm",
				Number: 0,
			},
			expectedUrl: "/supervision/deputies/firm/9952223",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.expectedUrl, test.firm.GetFirmURL())
		})
	}
}
