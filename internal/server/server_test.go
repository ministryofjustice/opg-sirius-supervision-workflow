package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"io"
	"time"
)

type mockTemplate struct {
	count    int
	lastVars interface{}
	lastW    io.Writer
	error    error
}

func (m *mockTemplate) Execute(w io.Writer, vars any) error {
	m.count += 1
	m.lastVars = vars
	m.lastW = w
	return m.error
}

type mockApiClient struct {
	error              error
	CurrentUserDetails model.Assignee
	TeamsForSelection  []model.Team
}

func (m mockApiClient) GetCurrentUserDetails(context sirius.Context) (model.Assignee, error) {
	return m.CurrentUserDetails, m.error
}

func (m mockApiClient) GetTeamsForSelection(context sirius.Context) ([]model.Team, error) {
	return m.TeamsForSelection, m.error
}

func (m mockApiClient) GetTaskTypes(context sirius.Context, strings []string) ([]model.TaskType, error) {
	panic("implement me")
}

func (m mockApiClient) GetTaskList(context sirius.Context, i int, i2 int, team model.Team, strings []string, types []model.TaskType, strings2 []string, time *time.Time, time2 *time.Time) (sirius.TaskList, error) {
	panic("implement me")
}

func (m mockApiClient) AssignTasksToCaseManager(context sirius.Context, i int, strings []string, s string) (string, error) {
	panic("implement me")
}

func (m mockApiClient) GetClientList(context sirius.Context, team model.Team, perPage int, page int) (sirius.ClientList, error) {
	panic("implement me")
}
