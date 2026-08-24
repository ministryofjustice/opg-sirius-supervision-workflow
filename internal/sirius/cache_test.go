package sirius

import (
	"testing"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCaches_updateAndGetTeams(t *testing.T) {
	caches := newCaches()
	expected := []model.Team{
		{Id: 1, Name: "Team One"},
		{Id: 2, Name: "Team Two"},
	}

	teams, found := caches.getTeams()
	assert.False(t, found)
	assert.Nil(t, teams)

	caches.updateTeams(expected)

	teams, found = caches.getTeams()
	assert.True(t, found)
	assert.Equal(t, expected, teams)
}

func TestCaches_updateAndGetTaskTypes(t *testing.T) {
	caches := newCaches()
	expected := []model.TaskType{
		{Handle: "A", Incomplete: "Task A"},
		{Handle: "B", Incomplete: "Task B"},
	}
	key := "supervision"

	taskTypes, found := caches.getTaskTypes(key)
	assert.False(t, found)
	assert.Nil(t, taskTypes)

	caches.updateTaskTypes(key, expected)

	taskTypes, found = caches.getTaskTypes(key)
	assert.True(t, found)
	assert.Equal(t, expected, taskTypes)
}
