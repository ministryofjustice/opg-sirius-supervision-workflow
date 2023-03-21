package sirius

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAppliedFiltersSingleTaskFilterSelectedReturned(t *testing.T) {
	apiTaskTypes := []ApiTaskTypes{
		{
			Handle:     "CWGN",
			Incomplete: "Casework - General",
			Complete:   "Casework - General",
			User:       true,
			Category:   "supervision",
			IsSelected: true,
		},
		{
			Handle:     "ORAL",
			Incomplete: "Order - Allocate to team",
			Complete:   "Order - Allocate to team",
			User:       true,
			Category:   "supervision",
			IsSelected: false,
		},
	}

	teamCollection := []ReturnedTeamCollection{
		{
			Id:        12,
			Name:      "Lay Team 1 - (Supervision)",
			Type:      "Supervision",
			TypeLabel: "Only",
			Members: []TeamMember{
				{
					ID:   1,
					Name: "Test One",
				},
				{
					ID:   2,
					Name: "Test Two",
				},
			},
			Selector: "12",
		},
		{
			Id:        13,
			Name:      "Allocations Team",
			Type:      "Supervision",
			TypeLabel: "Only",
			Selector:  "13",
		},
	}

	var selectedAssignees []string
	var selectedUnassigned string

	expectedFilter := []string{
		"Casework - General",
	}

	appliedFilters := GetAppliedFilters(teamCollection[0], selectedAssignees, selectedUnassigned, apiTaskTypes)

	assert.Equal(t, expectedFilter, appliedFilters)
	assert.Equal(t, 1, len(appliedFilters))
}

func TestGetAppliedFiltersMultipleTaskFilterSelectedReturned(t *testing.T) {
	apiTaskTypes := []ApiTaskTypes{
		{
			Handle:     "CWGN",
			Incomplete: "Casework - General",
			Complete:   "Casework - General",
			User:       true,
			Category:   "supervision",
			IsSelected: true,
		},
		{
			Handle:     "ORAL",
			Incomplete: "Order - Allocate to team",
			Complete:   "Order - Allocate to team",
			User:       true,
			Category:   "supervision",
			IsSelected: true,
		},
	}

	teamCollection := []ReturnedTeamCollection{
		{
			Id:        12,
			Name:      "Lay Team 1 - (Supervision)",
			Type:      "Supervision",
			TypeLabel: "Only",
			Selector:  "12",
		},
		{
			Id:        13,
			Name:      "Allocations Team",
			Type:      "Supervision",
			TypeLabel: "Only",
			Selector:  "13",
		},
	}

	var selectedAssignees []string
	var selectedUnassigned string

	expectedFilter := []string{
		"Casework - General",
		"Order - Allocate to team",
	}

	appliedFilters := GetAppliedFilters(teamCollection[0], selectedAssignees, selectedUnassigned, apiTaskTypes)

	assert.Equal(t, expectedFilter, appliedFilters)
	assert.Equal(t, 2, len(appliedFilters))
}

func TestGetAppliedFiltersSingleTaskSingleTeamMemberFilterSelectedReturned(t *testing.T) {
	apiTaskTypes := []ApiTaskTypes{
		{
			Handle:     "CWGN",
			Incomplete: "Casework - General",
			Complete:   "Casework - General",
			User:       true,
			Category:   "supervision",
			IsSelected: false,
		},
		{
			Handle:     "ORAL",
			Incomplete: "Order - Allocate to team",
			Complete:   "Order - Allocate to team",
			User:       true,
			Category:   "supervision",
			IsSelected: true,
		},
	}

	teamCollection := []ReturnedTeamCollection{
		{
			Id:        12,
			Name:      "Supervision Team 1",
			Type:      "Supervision",
			TypeLabel: "Only",
			Selector:  "12",
		},
		{
			Id:        13,
			Name:      "Allocations Team",
			Type:      "Supervision",
			TypeLabel: "Only",
			Members: []TeamMember{
				{
					ID:   1,
					Name: "Test One",
				},
				{
					ID:   2,
					Name: "Test Two",
				},
			},
			Selector: "13",
		},
	}

	selectedAssignees := []string{"2"}
	var selectedUnassigned string

	expectedFilter := []string{
		"Order - Allocate to team",
		"Test Two",
	}

	appliedFilters := GetAppliedFilters(teamCollection[1], selectedAssignees, selectedUnassigned, apiTaskTypes)

	assert.Equal(t, expectedFilter, appliedFilters)
	assert.Equal(t, 2, len(appliedFilters))
}

func TestGetAppliedFiltersMultipleTasksSingleTeamMemberAndUnassignedFilterSelectedReturned(t *testing.T) {
	apiTaskTypes := []ApiTaskTypes{
		{
			Handle:     "CWGN",
			Incomplete: "Casework - General",
			Complete:   "Casework - General",
			User:       true,
			Category:   "supervision",
			IsSelected: true,
		},
		{
			Handle:     "ORAL",
			Incomplete: "Order - Allocate to team",
			Complete:   "Order - Allocate to team",
			User:       true,
			Category:   "supervision",
			IsSelected: true,
		},
	}

	teamCollection := []ReturnedTeamCollection{
		{
			Id:        12,
			Name:      "Supervision Team 1",
			Type:      "Supervision",
			TypeLabel: "Only",
			Selector:  "12",
		},
		{
			Id:        13,
			Name:      "Allocations Team",
			Type:      "Supervision",
			TypeLabel: "Only",
			Members: []TeamMember{
				{
					ID:   1,
					Name: "Test One",
				},
				{
					ID:   2,
					Name: "Test Two",
				},
			},
			Selector: "13",
		},
	}

	selectedAssignees := []string{"1"}
	selectedUnassigned := teamCollection[1].Selector

	expectedFilter := []string{
		"Casework - General",
		"Order - Allocate to team",
		"Allocations Team",
		"Test One",
	}

	appliedFilters := GetAppliedFilters(teamCollection[1], selectedAssignees, selectedUnassigned, apiTaskTypes)

	assert.Equal(t, expectedFilter, appliedFilters)
	assert.Equal(t, 4, len(appliedFilters))
}

func TestGetAppliedFiltersNoFiltersSelectedReturned(t *testing.T) {
	apiTaskTypes := []ApiTaskTypes{
		{
			Handle:     "CWGN",
			Incomplete: "Casework - General",
			Complete:   "Casework - General",
			User:       true,
			Category:   "supervision",
			IsSelected: false,
		},
		{
			Handle:     "ORAL",
			Incomplete: "Order - Allocate to team",
			Complete:   "Order - Allocate to team",
			User:       true,
			Category:   "supervision",
			IsSelected: false,
		},
	}

	teamCollection := []ReturnedTeamCollection{
		{
			Id:        12,
			Name:      "Supervision Team 1",
			Type:      "Supervision",
			TypeLabel: "Only",
			Selector:  "12",
		},
		{
			Id:        13,
			Name:      "Allocations Team",
			Type:      "Supervision",
			TypeLabel: "Only",
			Members: []TeamMember{
				{
					ID:   1,
					Name: "Test One",
				},
				{
					ID:   2,
					Name: "Test Two",
				},
			},
			Selector: "13",
		},
	}

	var selectedAssignees []string
	var selectedUnassigned string
	var expectedFilter []string

	appliedFilters := GetAppliedFilters(teamCollection[1], selectedAssignees, selectedUnassigned, apiTaskTypes)

	assert.Equal(t, expectedFilter, appliedFilters)
	assert.Equal(t, 0, len(appliedFilters))
}
