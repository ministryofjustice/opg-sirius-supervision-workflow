package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockWorkflowClient struct {
	count        int
	lastCtx      sirius.Context
	err          error
	userData     sirius.UserDetails
	taskTypeData sirius.WholeTaskList
	taskListData sirius.TaskList
}

func (m *mockWorkflowClient) SiriusUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userData, m.err
}

func (m *mockWorkflowClient) GetTaskDetails(ctx sirius.Context) (sirius.WholeTaskList, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.taskTypeData, m.err
}

func (m *mockWorkflowClient) GetTaskList(ctx sirius.Context) (sirius.TaskList, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.taskListData, m.err
}

func TestGetUserDetails(t *testing.T) {
	assert := assert.New(t)

	mockUserDetailsData := sirius.UserDetails{
		ID:        123,
		Firstname: "John",
		Surname:   "Doe",
	}

	mockTaskTypeData := sirius.WholeTaskList{
		AllTaskList: sirius.ApiTaskTypes{
			Handle:     "CDFC",
			Incomplete: "Correspondence - Review failed draft",
			Category:   "supervision",
			Complete:   "Correspondence - Reviewed draft failure",
			User:       true,
		},
	}

	mockTaskListData := sirius.TaskList{
		AllTaskList: []sirius.ApiTask{
			{
				ApiTaskAssignee: sirius.AssigneeDetails{
					AssigneeDetailsDisplayName: "Assignee Duke Clive Henry Hetley Junior Jones",
				},
				ApiTaskType:    "Case work - General",
				ApiTaskDueDate: "01/02/2021",
				ApiTaskCaseItems: []sirius.CaseItemsDetails{
					{
						CaseItemClient: sirius.ClientDetails{
							ClientDetailsCaseRecNumber: "caseRecNumber",
							ClientDetailsFirstName:     "Client Alexander Zacchaeus",
							ClientDetailsId:            3333,
							ClientDetailsSupervisionCaseOwner: sirius.SupervisionCaseOwnerDetail{
								SupervisionCaseOwnerName: "Supervision - Team - Name",
							},
							ClientDetailsSurname: "Client Wolfeschlegelsteinhausenbergerdorff",
						},
					},
				},
			},
		},
	}

	client := &mockWorkflowClient{userData: mockUserDetailsData, taskTypeData: mockTaskTypeData, taskListData: mockTaskListData}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := loggingInfoForWorflow(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(3, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(userDetailsVars{
		Path:      "/path",
		ID:        123,
		Firstname: "John",
		Surname:   "Doe",
		TaskList: sirius.TaskList{
			AllTaskList: []sirius.ApiTask{
				{
					ApiTaskAssignee: sirius.AssigneeDetails{
						AssigneeDetailsDisplayName: "Assignee Duke Clive Henry Hetley Junior Jones",
					},
					ApiTaskCaseItems: []sirius.CaseItemsDetails{
						{
							CaseItemClient: sirius.ClientDetails{
								ClientDetailsCaseRecNumber: "caseRecNumber",
								ClientDetailsFirstName:     "Client Alexander Zacchaeus",
								ClientDetailsId:            3333,
								ClientDetailsSupervisionCaseOwner: sirius.SupervisionCaseOwnerDetail{
									SupervisionCaseOwnerName: "Supervision - Team - Name",
								},
								ClientDetailsSurname: "Client Wolfeschlegelsteinhausenbergerdorff",
							},
						},
					},
					ApiTaskDueDate: "01/02/2021",
					ApiTaskType:    "Case work - General",
				},
			},
		},
		LoadTasks: sirius.WholeTaskList{
			AllTaskList: sirius.ApiTaskTypes{
				Handle:     "CDFC",
				Incomplete: "Correspondence - Review failed draft",
				Category:   "supervision",
				Complete:   "Correspondence - Reviewed draft failure",
				User:       true,
			},
		},
	}, template.lastVars)

}

func TestWorkflowUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockWorkflowClient{err: sirius.ErrUnauthorized}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := loggingInfoForWorflow(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal(sirius.ErrUnauthorized, err)

	assert.Equal(0, template.count)
}

func TestWorkflowSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockWorkflowClient{err: errors.New("err")}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := loggingInfoForWorflow(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal("err", err.Error())

	assert.Equal(0, template.count)
}

func TestPostWorkflow(t *testing.T) {
	assert := assert.New(t)
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)

	handler := loggingInfoForWorflow(nil, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, template.count)
}
