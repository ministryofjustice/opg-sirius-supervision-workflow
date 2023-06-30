package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTeam_GetAssigneesForFilter(t *testing.T) {
	team := Team{
		Members: []Assignee{
			{Id: 1, Name: "B"},
			{Id: 2, Name: "A"},
		},
		Teams: []Team{
			{
				Members: []Assignee{
					{Id: 4, Name: "D"},
					{Id: 2, Name: "A"},
					{Id: 3, Name: "C"},
				},
			},
			{
				Members: []Assignee{
					{Id: 3, Name: "C"},
				},
			},
		},
	}

	expected := []Assignee{
		{Id: 2, Name: "A"},
		{Id: 1, Name: "B"},
		{Id: 3, Name: "C"},
		{Id: 4, Name: "D"},
	}

	assert.Equal(t, expected, team.GetAssigneesForFilter())
}

func TestTeam_HasTeam(t *testing.T) {
	team := Team{
		Id: 10,
		Teams: []Team{
			{Id: 12},
			{Id: 13},
		},
	}

	assert.Truef(t, team.HasTeam(10), "Parent team ID 10 not found")
	assert.Truef(t, team.HasTeam(12), "Check team ID 12 not found")
	assert.Truef(t, team.HasTeam(13), "Child team ID 13 not found")
	assert.False(t, team.HasTeam(11), "Child team ID 11 should not exist")
}

func TestTeam_IsLay(t *testing.T) {
	assert.True(t, Team{Type: "LAY"}.IsLay())
	assert.True(t, Team{Type: "", Selector: "lay-team"}.IsLay())
	assert.False(t, Team{Type: "NOT LAY"}.IsLay())
}

func TestTeam_IsPro(t *testing.T) {
	assert.True(t, Team{Type: "PRO"}.IsPro())
	assert.True(t, Team{Type: "", Selector: "pro-team"}.IsPro())
	assert.False(t, Team{Type: "NOT PRO"}.IsPro())
}

func TestTeam_IsLayNewOrdersTeam(t *testing.T) {
	assert.False(t, Team{Name: "Random team"}.IsLayNewOrdersTeam())
	assert.True(t, Team{Name: "Lay Team - New Deputy Orders"}.IsLayNewOrdersTeam())
}
