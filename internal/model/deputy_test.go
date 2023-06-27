package model

import (
	"github.com/stretchr/testify/assert"
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
