package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"io"
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

func (m mockApiClient) ReassignClients(context sirius.Context, params sirius.ReassignClientsParams) (string, error) {
	return "", nil
}

func (m mockApiClient) GetCurrentUserDetails(context sirius.Context) (model.Assignee, error) {
	return m.CurrentUserDetails, m.error
}

func (m mockApiClient) GetTeamsForSelection(context sirius.Context, teamType []string) ([]model.Team, error) {
	return m.TeamsForSelection, m.error
}

func (m mockApiClient) GetTaskTypes(context sirius.Context, params sirius.TaskTypesParams) ([]model.TaskType, error) {
	panic("implement me")
}

func (m mockApiClient) GetTaskList(context sirius.Context, params sirius.TaskListParams) (sirius.TaskList, error) {
	panic("implement me")
}

func (m mockApiClient) ReassignTasks(context sirius.Context, params sirius.ReassignTasksParams) (string, error) {
	panic("implement me")
}

func (m mockApiClient) GetClientList(context sirius.Context, params sirius.ClientListParams) (sirius.ClientList, error) {
	panic("implement me")
}

func (m mockApiClient) GetDeputyList(context sirius.Context, params sirius.DeputyListParams) (sirius.DeputyList, error) {
	panic("implement me")
}

func (m mockApiClient) ReassignDeputies(context sirius.Context, params sirius.ReassignDeputiesParams) (string, error) {
	return "", nil
}
