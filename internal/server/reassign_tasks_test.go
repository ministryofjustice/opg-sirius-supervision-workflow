package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReassignTasksClient struct {
	mock.Mock
}

func (m *mockReassignTasksClient) ReassignTasks(ctx sirius.Context, params sirius.ReassignTasksParams) (string, error) {
	args := m.Called(ctx)
	return args.Get(0).(string), args.Error(1)
}

func TestReassignTasks_ClientTasks(t *testing.T) {
	client := &mockReassignTasksClient{}

	expectedParams := sirius.ReassignTasksParams{
		AssignTeam: "10",
		AssignCM:   "20",
		TaskIds:    []string{"1", "2"},
		IsPriority: "true",
	}

	client.On("ReassignTasks", mock.Anything).Return("reassign successful", nil)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/client-tasks?team=19&page=1&per-page=25&task-type=CDFC&task-type=ORAL", nil)
	r.PostForm = url.Values{
		"assignTeam":     {expectedParams.AssignTeam},
		"assignCM":       {expectedParams.AssignCM},
		"selected-tasks": expectedParams.TaskIds,
		"priority":       {expectedParams.IsPriority},
	}

	app := WorkflowVars{
		Path:         "/client-tasks?team=19&page=1&per-page=25&task-type=CDFC&task-type=ORAL",
		SelectedTeam: model.Team{Type: "LAY", Selector: "19", Id: 19},
		MyDetails: model.Assignee{
			Teams: []model.Team{
				{
					Id:   99,
					Name: "my-team",
				},
			},
			Roles: []string{"Case Manager"},
		},
	}
	err := reassignTasks(client)(app, w, r)

	assert.Equal(t, Redirect{
		Path:           "/client-tasks?team=19&page=1&per-page=25&task-type=CDFC&task-type=ORAL",
		SuccessMessage: "reassign successful",
	}, err)
}

func TestReassignTasks_DeputyTasks(t *testing.T) {
	client := &mockReassignTasksClient{}

	client.On("ReassignTasks", mock.Anything).Return("reassign success", nil)

	expectedParams := sirius.ReassignTasksParams{
		AssignTeam: "10",
		AssignCM:   "20",
		TaskIds:    []string{"1", "2"},
		IsPriority: "true",
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/deputy-tasks?team=1&page=2&per-page=25&order-by=deputy&sort=asc", nil)
	r.PostForm = url.Values{
		"assignTeam":     {expectedParams.AssignTeam},
		"assignCM":       {expectedParams.AssignCM},
		"selected-tasks": expectedParams.TaskIds,
		"priority":       {expectedParams.IsPriority},
	}

	workflowVars := WorkflowVars{
		MyDetails: model.Assignee{
			Id: 123,
		},
		Path:         "deputy-tasks",
		SelectedTeam: model.Team{Type: "PRO", Selector: "1"},
	}
	err := reassignTasks(client)(workflowVars, w, r)
	assert.Equal(t, Redirect{
		Path:           "/deputy-tasks?team=1&page=2&per-page=25&order-by=deputy&sort=asc",
		SuccessMessage: "reassign success",
	}, err)
}
