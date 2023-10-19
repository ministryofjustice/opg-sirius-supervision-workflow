package model

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestAssurance_IsPDR(t *testing.T) {
	assert.False(t, Assurance{}.IsPDR())
	assert.False(t, Assurance{Type: RefData{Handle: "notPDR"}}.IsPDR())
	assert.True(t, Assurance{Type: RefData{Handle: "PDR"}}.IsPDR())
}

func TestAssurance_GetRAGRating(t *testing.T) {
	tests := []struct {
		assurance Assurance
		want      RAGRating
	}{
		{
			assurance: Assurance{},
			want:      RAGRating{},
		},
		{
			assurance: Assurance{ReportMarkedAs: RefData{Handle: "RED"}},
			want:      RAGRating{Name: "High risk", Colour: "red"},
		},
		{
			assurance: Assurance{ReportMarkedAs: RefData{Handle: "AMBER"}},
			want:      RAGRating{Name: "Medium risk", Colour: "orange"},
		},
		{
			assurance: Assurance{ReportMarkedAs: RefData{Handle: "GREEN"}},
			want:      RAGRating{Name: "Low risk", Colour: "green"},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.want, test.assurance.GetRAGRating())
		})
	}
}
