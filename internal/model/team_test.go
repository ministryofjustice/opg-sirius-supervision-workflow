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

func TestTeam_IsPA(t *testing.T) {
	assert.True(t, Team{Type: "PA"}.IsPA())
	assert.False(t, Team{Type: "", Selector: "pa-team"}.IsPA())
	assert.False(t, Team{Type: "NOT PA"}.IsPA())
}

func TestTeam_IsLayNewOrdersTeam(t *testing.T) {
	assert.False(t, Team{Name: "Random team"}.IsLayNewOrdersTeam())
	assert.True(t, Team{Name: "Lay Team - New Deputy Orders"}.IsLayNewOrdersTeam())
}

func TestTeam_IsHW(t *testing.T) {
	assert.True(t, Team{Type: "HW"}.IsHW())
	assert.False(t, Team{Type: "NOT HW"}.IsHW())
}

func TestTeam_IsClosedCasesTeam(t *testing.T) {
	assert.False(t, Team{Name: "Random team"}.IsClosedCases())
	assert.True(t, Team{Name: "supervision closed cases"}.IsClosedCases())
	assert.True(t, Team{Name: "Supervision closed cases"}.IsClosedCases())
	assert.True(t, Team{Name: "Supervision Closed Cases"}.IsClosedCases())
}

func TestGetUnassignedCount(t *testing.T) {
	tests := []struct {
		testname     string
		selectedTeam Team
		want         string
		url          string
	}{
		{
			testname:     "Returns count for team",
			selectedTeam: Team{Id: 10},
			want:         "(11)",
			url:          "test",
		},
		{
			testname:     "Returns null count for team",
			selectedTeam: Team{Id: 11},
			want:         "(0)",
			url:          "test",
		},
		{
			testname:     "Returns null if team not in list",
			selectedTeam: Team{Id: 22},
			want:         "(0)",
			url:          "test",
		},
		//{
		//	testname:     "Returns empty string if deputy tasks page",
		//	selectedTeam: Team{Id: 22},
		//	want:         "",
		//	url:          "deputy-tasks",
		//},
	}
	for _, test := range tests {
		selectedAssignees := []AssigneeAndCount{
			{AssigneeId: 10, Count: 11},
			{AssigneeId: 11, Count: 0},
			{AssigneeId: 12, Count: 1},
		}

		t.Run(test.testname, func(t *testing.T) {
			assert.Equal(t, test.selectedTeam.GetUnassignedCount(selectedAssignees, test.url), test.want)
		})
	}
}

func TestGetCountForATeam(t *testing.T) {
	tests := []struct {
		testname string
		teamId   int
		want     int
	}{
		{
			testname: "Returns count for team",
			teamId:   10,
			want:     20,
		},
		{
			testname: "Returns null count for team",
			teamId:   11,
			want:     0,
		},
		{
			testname: "Returns null if team not in list",
			teamId:   13,
			want:     0,
		},
	}
	for _, test := range tests {
		selectedAssignees := []AssigneeAndCount{
			{AssigneeId: 10, Count: 20},
			{AssigneeId: 11, Count: 0},
			{AssigneeId: 12, Count: 1},
		}

		t.Run(test.testname, func(t *testing.T) {
			myTeam := Team{Id: 99}
			assert.Equal(t, myTeam.GetCountForATeam(selectedAssignees, test.teamId), test.want)
		})
	}
}

func TestGetMultiTeamUnassignedCount(t *testing.T) {
	tests := []struct {
		testname string
		teams    []Team
		want     int
		url      string
	}{
		{
			testname: "Returns count for teams",
			teams: []Team{
				{Id: 20},
				{Id: 30},
			},
			want: 21,
			url:  "test",
		},
		{
			testname: "Returns count even when some team ids dont exist",
			teams: []Team{
				{Id: 99},
				{Id: 30},
				{Id: 98},
			},
			want: 1,
			url:  "test",
		},
		{
			testname: "Returns null if no teams found",
			teams: []Team{
				{Id: 97},
				{Id: 98},
				{Id: 99},
			},
			want: 0,
		},
	}
	for _, test := range tests {
		selectedAssignees := []AssigneeAndCount{
			{AssigneeId: 10, Count: 20},
			{AssigneeId: 11, Count: 0},
			{AssigneeId: 12, Count: 1},
			{AssigneeId: 20, Count: 20},
			{AssigneeId: 21, Count: 0},
			{AssigneeId: 30, Count: 1},
		}

		t.Run(test.testname, func(t *testing.T) {
			myTeam := Team{Id: 99}
			myTeam.Teams = test.teams
			assert.Equal(t, myTeam.GetMultiTeamUnassignedCount(selectedAssignees), test.want)
		})
	}
}
