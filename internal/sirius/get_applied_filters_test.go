package sirius

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAppliedFiltersSingleTaskFilterSelectedReturned(t *testing.T) {

	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

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
			Id:             12,
			Name:           "Lay Team 1 - (Supervision)",
			Type:           "Supervision",
			TypeLabel:      "Only",
			IsTeamSelected: false,
		},
		{
			Id:             13,
			Name:           "Allocations Team",
			Type:           "Supervision",
			TypeLabel:      "Only",
			IsTeamSelected: false,
		},
	}

	assigneeTeam := AssigneesTeam{
		Id:   12,
		Name: "Lay Team 1 - (Supervision)",
		Members: []AssigneeTeamMembers{
			{
				TeamMembersId:          1,
				TeamMembersName:        "Test One",
				TeamMembersDisplayName: "Test One",
				IsSelected:             false,
			},
			{
				TeamMembersId:          2,
				TeamMembersName:        "Test Two",
				TeamMembersDisplayName: "Test Two",
				IsSelected:             false,
			},
		},
	}

	expectedFilter := []string{
		"Casework - General",
	}

	appliedFilters := client.GetAppliedFilters(12, apiTaskTypes, teamCollection, assigneeTeam)

	assert.Equal(t, expectedFilter, appliedFilters)
	assert.Equal(t, 1, len(appliedFilters))
}

func TestGetAppliedFiltersMultipleTaskFilterSelectedReturned(t *testing.T) {

	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

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
			Id:             12,
			Name:           "Lay Team 1 - (Supervision)",
			Type:           "Supervision",
			TypeLabel:      "Only",
			IsTeamSelected: false,
		},
		{
			Id:             13,
			Name:           "Allocations Team",
			Type:           "Supervision",
			TypeLabel:      "Only",
			IsTeamSelected: false,
		},
	}

	assigneeTeam := AssigneesTeam{}

	expectedFilter := []string{
		"Casework - General",
		"Order - Allocate to team",
	}

	appliedFilters := client.GetAppliedFilters(12, apiTaskTypes, teamCollection, assigneeTeam)

	assert.Equal(t, expectedFilter, appliedFilters)
	assert.Equal(t, 2, len(appliedFilters))
}

func TestGetAppliedFiltersSingleTaskSingleTeamFilterSelectedReturned(t *testing.T) {

	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

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
			Id:             12,
			Name:           "Supervision Team 1",
			Type:           "Supervision",
			TypeLabel:      "Only",
			IsTeamSelected: false,
		},
		{
			Id:             13,
			Name:           "Allocations Team",
			Type:           "Supervision",
			TypeLabel:      "Only",
			IsTeamSelected: true,
		},
	}

	assigneeTeam := AssigneesTeam{
		Id:   13,
		Name: "Lay Team 1 - (Supervision)",
		Members: []AssigneeTeamMembers{
			{
				TeamMembersId:          1,
				TeamMembersName:        "Test One",
				TeamMembersDisplayName: "Test One",
				IsSelected:             false,
			},
			{
				TeamMembersId:          2,
				TeamMembersName:        "Test Two",
				TeamMembersDisplayName: "Test Two",
				IsSelected:             false,
			},
		},
	}

	expectedFilter := []string{
		"Order - Allocate to team",
		"Allocations Team",
	}

	appliedFilters := client.GetAppliedFilters(13, apiTaskTypes, teamCollection, assigneeTeam)

	assert.Equal(t, expectedFilter, appliedFilters)
	assert.Equal(t, 2, len(appliedFilters))
}

func TestGetAppliedFiltersSingleTaskSingleTeamSingleTeamMemberFilterSelectedReturned(t *testing.T) {

	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

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
			Id:             12,
			Name:           "Supervision Team 1",
			Type:           "Supervision",
			TypeLabel:      "Only",
			IsTeamSelected: false,
		},
		{
			Id:             13,
			Name:           "Allocations Team",
			Type:           "Supervision",
			TypeLabel:      "Only",
			IsTeamSelected: true,
		},
	}

	assigneeTeam := AssigneesTeam{
		Id:   13,
		Name: "Allocations Team",
		Members: []AssigneeTeamMembers{
			{
				TeamMembersId:          1,
				TeamMembersName:        "Test One",
				TeamMembersDisplayName: "Test One",
				IsSelected:             false,
			},
			{
				TeamMembersId:          2,
				TeamMembersName:        "Test Two",
				TeamMembersDisplayName: "Test Two",
				IsSelected:             true,
			},
		},
	}

	expectedFilter := []string{
		"Order - Allocate to team",
		"Allocations Team",
		"Test Two",
	}

	appliedFilters := client.GetAppliedFilters(13, apiTaskTypes, teamCollection, assigneeTeam)

	assert.Equal(t, expectedFilter, appliedFilters)
	assert.Equal(t, 3, len(appliedFilters))
}

func TestGetAppliedFiltersMultipleTasksMultipleTeamsSingleTeamMemberFilterSelectedReturned(t *testing.T) {

	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

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
			Id:             12,
			Name:           "Supervision Team 1",
			Type:           "Supervision",
			TypeLabel:      "Only",
			IsTeamSelected: true,
		},
		{
			Id:             13,
			Name:           "Allocations Team",
			Type:           "Supervision",
			TypeLabel:      "Only",
			IsTeamSelected: true,
		},
	}

	assigneeTeam := AssigneesTeam{
		Id:   13,
		Name: "Allocations Team",
		Members: []AssigneeTeamMembers{
			{
				TeamMembersId:          1,
				TeamMembersName:        "Test One",
				TeamMembersDisplayName: "Test One",
				IsSelected:             true,
			},
			{
				TeamMembersId:          2,
				TeamMembersName:        "Test Two",
				TeamMembersDisplayName: "Test Two",
				IsSelected:             false,
			},
		},
	}

	expectedFilter := []string{
		"Casework - General",
		"Order - Allocate to team",
		"Allocations Team",
		"Test One",
	}

	appliedFilters := client.GetAppliedFilters(13, apiTaskTypes, teamCollection, assigneeTeam)

	assert.Equal(t, expectedFilter, appliedFilters)
	assert.Equal(t, 4, len(appliedFilters))
}

func TestGetAppliedFiltersNoFiltersSelectedReturned(t *testing.T) {

	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

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
			Id:             12,
			Name:           "Supervision Team 1",
			Type:           "Supervision",
			TypeLabel:      "Only",
			IsTeamSelected: false,
		},
		{
			Id:             13,
			Name:           "Allocations Team",
			Type:           "Supervision",
			TypeLabel:      "Only",
			IsTeamSelected: false,
		},
	}

	assigneeTeam := AssigneesTeam{
		Id:   13,
		Name: "Allocations Team",
		Members: []AssigneeTeamMembers{
			{
				TeamMembersId:          1,
				TeamMembersName:        "Test One",
				TeamMembersDisplayName: "Test One",
				IsSelected:             false,
			},
			{
				TeamMembersId:          2,
				TeamMembersName:        "Test Two",
				TeamMembersDisplayName: "Test Two",
				IsSelected:             false,
			},
		},
	}

	expectedFilter := []string(nil)

	appliedFilters := client.GetAppliedFilters(13, apiTaskTypes, teamCollection, assigneeTeam)

	assert.Equal(t, expectedFilter, appliedFilters)
	assert.Equal(t, 0, len(appliedFilters))
}
