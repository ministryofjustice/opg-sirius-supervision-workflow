package sirius

import (
	"time"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/patrickmn/go-cache"
)

const (
	teamsDuration     = 5 * time.Minute
	teamsKey          = "teams"
	taskTypesDuration = 30 * time.Minute
)

type Caches struct {
	teams     *cache.Cache
	taskTypes *cache.Cache
}

func newCaches() *Caches {
	teams := cache.New(teamsDuration, teamsDuration)
	_ = teams.Add("refresh", false, teamsDuration)
	taskTypes := cache.New(taskTypesDuration, taskTypesDuration)
	_ = taskTypes.Add("refresh", false, taskTypesDuration)

	return &Caches{
		teams:     teams,
		taskTypes: taskTypes,
	}
}

func (c Caches) updateTeams(teams []model.Team) {
	_ = c.teams.Add(teamsKey, &teams, teamsDuration)
}

func (c Caches) getTeams() ([]model.Team, bool) {
	if x, found := c.teams.Get(teamsKey); found {
		return *(x.(*[]model.Team)), true
	}
	return nil, false
}

func (c Caches) updateTaskTypes(key string, taskTypes []model.TaskType) {
	_ = c.taskTypes.Add(key, &taskTypes, taskTypesDuration)
}

func (c Caches) getTaskTypes(key string) ([]model.TaskType, bool) {
	if x, found := c.taskTypes.Get(key); found {
		return *(x.(*[]model.TaskType)), true
	}
	return nil, false
}
