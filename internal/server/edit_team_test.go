package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockEditTeamClient struct {
	team struct {
		count   int
		lastCtx sirius.Context
		lastID  int
		data    sirius.Team
		err     error
	}

	teamTypes struct {
		count   int
		lastCtx sirius.Context
		data    []sirius.RefDataTeamType
		err     error
	}

	editTeam struct {
		count    int
		lastCtx  sirius.Context
		lastTeam sirius.Team
		err      error
	}
}

func (m *mockEditTeamClient) Team(ctx sirius.Context, id int) (sirius.Team, error) {
	m.team.count += 1
	m.team.lastCtx = ctx
	m.team.lastID = id

	return m.team.data, m.team.err
}

func (m *mockEditTeamClient) TeamTypes(ctx sirius.Context) ([]sirius.RefDataTeamType, error) {
	m.teamTypes.count += 1
	m.teamTypes.lastCtx = ctx

	return m.teamTypes.data, m.teamTypes.err
}

func (m *mockEditTeamClient) EditTeam(ctx sirius.Context, team sirius.Team) error {
	m.editTeam.count += 1
	m.editTeam.lastCtx = ctx
	m.editTeam.lastTeam = team

	return m.editTeam.err
}

func (m *mockEditTeamClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{
		"team": sirius.PermissionGroup{Permissions: []string{"put", "post"}},
	}
}

func TestGetEditTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditTeamClient{}
	client.team.data = sirius.Team{DisplayName: "Complaints team"}
	client.teamTypes.data = []sirius.RefDataTeamType{
		{
			Handle: "TEST",
			Label:  "Test type",
		},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/edit/123", nil)

	err := editTeam(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.team.count)
	assert.Equal(123, client.team.lastID)

	assert.Equal(1, client.teamTypes.count)

	assert.Equal(0, client.editTeam.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editTeamVars{
		Path:            "/teams/edit/123",
		Team:            client.team.data,
		TeamTypeOptions: client.teamTypes.data,
		CanEditTeamType: true,
	}, template.lastVars)
}

func TestGetEditTeamNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := editTeam(nil, nil)(sirius.PermissionSet{}, w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
}

func TestGetEditTeamWithoutTypeEditPermission(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditTeamClient{}
	client.team.data = sirius.Team{DisplayName: "Complaints team"}
	client.teamTypes.data = []sirius.RefDataTeamType{
		{
			Handle: "TEST",
			Label:  "Test type",
		},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/edit/123", nil)

	err := editTeam(client, template)(sirius.PermissionSet{
		"team": sirius.PermissionGroup{Permissions: []string{"put"}},
	}, w, r)
	assert.Nil(err)

	assert.Equal(editTeamVars{
		Path:            "/teams/edit/123",
		Team:            client.team.data,
		TeamTypeOptions: client.teamTypes.data,
	}, template.lastVars)
}

func TestGetEditTeamWithDeletePermission(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditTeamClient{}
	client.team.data = sirius.Team{DisplayName: "Complaints team"}
	client.teamTypes.data = []sirius.RefDataTeamType{
		{
			Handle: "TEST",
			Label:  "Test type",
		},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/edit/123", nil)

	permissions := client.requiredPermissions()
	permissions["v1-teams"] = sirius.PermissionGroup{Permissions: []string{"delete"}}

	err := editTeam(client, template)(permissions, w, r)
	assert.Nil(err)

	assert.Equal(editTeamVars{
		Path:            "/teams/edit/123",
		Team:            client.team.data,
		TeamTypeOptions: client.teamTypes.data,
		CanEditTeamType: true,
		CanDeleteTeam:   true,
	}, template.lastVars)
}

func TestGetEditTeamBadPath(t *testing.T) {
	for name, path := range map[string]string{
		"empty":       "/teams/edit/",
		"non-numeric": "/teams/edit/hello",
		"suffixed":    "/teams/edit/123/no",
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockEditTeamClient{}
			client.team.data = sirius.Team{DisplayName: "Complaints team"}
			template := &mockTemplate{}

			r, _ := http.NewRequest("GET", path, nil)

			err := editTeam(client, template)(client.requiredPermissions(), nil, r)

			assert.Equal(StatusError(http.StatusNotFound), err)

			assert.Equal(0, client.team.count)
			assert.Equal(0, client.editTeam.count)
			assert.Equal(0, template.count)
		})
	}
}

func TestPostEditTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditTeamClient{}
	client.team.data = sirius.Team{
		ID:          123,
		DisplayName: "Complaints team",
		Type:        "COMPLAINTS",
		Email:       "complaint@opgtest.com",
		PhoneNumber: "01234",
	}
	client.teamTypes.data = []sirius.RefDataTeamType{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/edit/123", strings.NewReader("name=New+name&service=supervision&supervision-type=FINANCE&email=new@opgtest.com&phone=9876"))
	r.Header.Add("Content-type", "application/x-www-form-urlencoded")

	err := editTeam(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.team.count)
	assert.Equal(123, client.team.lastID)
	assert.Equal(1, client.editTeam.count)
	assert.Equal(123, client.editTeam.lastTeam.ID)
	assert.Equal("New name", client.editTeam.lastTeam.DisplayName)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editTeamVars{
		Path: "/teams/edit/123",
		Team: sirius.Team{
			ID:          123,
			DisplayName: "New name",
			Type:        "FINANCE",
			Email:       "new@opgtest.com",
			PhoneNumber: "9876",
		},
		TeamTypeOptions: client.teamTypes.data,
		CanEditTeamType: true,
		Success:         true,
	}, template.lastVars)
}

func TestPostEditLpaTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditTeamClient{}
	client.team.data = sirius.Team{
		ID:          123,
		DisplayName: "Complaints team",
		Type:        "COMPLAINTS",
		Email:       "complaint@opgtest.com",
		PhoneNumber: "01234",
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/edit/123", strings.NewReader("name=New+name&service=lpa&email=new@opgtest.com&phone=9876"))
	r.Header.Add("Content-type", "application/x-www-form-urlencoded")

	err := editTeam(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.team.count)
	assert.Equal(123, client.team.lastID)
	assert.Equal(1, client.editTeam.count)
	assert.Equal(123, client.editTeam.lastTeam.ID)
	assert.Equal("New name", client.editTeam.lastTeam.DisplayName)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editTeamVars{
		Path: "/teams/edit/123",
		Team: sirius.Team{
			ID:          123,
			DisplayName: "New name",
			Type:        "",
			Email:       "new@opgtest.com",
			PhoneNumber: "9876",
		},
		TeamTypeOptions: client.teamTypes.data,
		CanEditTeamType: true,
		Success:         true,
	}, template.lastVars)
}

func TestPostEditTeamWithoutPermission(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditTeamClient{}
	client.team.data = sirius.Team{
		ID:          123,
		DisplayName: "Complaints team",
		Type:        "COMPLAINTS",
		Email:       "complaint@opgtest.com",
		PhoneNumber: "01234",
	}
	client.teamTypes.data = []sirius.RefDataTeamType{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/edit/123", strings.NewReader("name=New+name&service=lpa&email=new@opgtest.com&phone=9876"))
	r.Header.Add("Content-type", "application/x-www-form-urlencoded")

	err := editTeam(client, template)(sirius.PermissionSet{
		"team": sirius.PermissionGroup{Permissions: []string{"put"}},
	}, w, r)
	assert.Nil(err)

	assert.Equal(1, client.editTeam.count)
	assert.Equal("COMPLAINTS", client.editTeam.lastTeam.Type)

	assert.Equal(editTeamVars{
		Path: "/teams/edit/123",
		Team: sirius.Team{
			ID:          123,
			DisplayName: "New name",
			Type:        "COMPLAINTS",
			Email:       "new@opgtest.com",
			PhoneNumber: "9876",
		},
		TeamTypeOptions: client.teamTypes.data,
		Success:         true,
	}, template.lastVars)
}

func TestPostEditTeamValidationError(t *testing.T) {
	assert := assert.New(t)

	validationErrors := sirius.ValidationErrors{
		"teamType": {
			"invalidTeamType": "Invalid team type",
		},
	}

	client := &mockEditTeamClient{}
	client.team.data = sirius.Team{
		ID:          123,
		DisplayName: "Complaints team",
		Type:        "COMPLAINTS",
		Email:       "complaint@opgtest.com",
		PhoneNumber: "01234",
	}
	client.editTeam.err = &sirius.ValidationError{
		Errors: validationErrors,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/edit/123", strings.NewReader("name=New+name&service=supervision&supervision-type=FINANCE&email=new@opgtest.com&phone=9876"))
	r.Header.Add("Content-type", "application/x-www-form-urlencoded")

	err := editTeam(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(http.StatusBadRequest, w.Result().StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editTeamVars{
		Path: "/teams/edit/123",
		Team: sirius.Team{
			ID:          123,
			DisplayName: "New name",
			Type:        "FINANCE",
			Email:       "new@opgtest.com",
			PhoneNumber: "9876",
		},
		TeamTypeOptions: client.teamTypes.data,
		CanEditTeamType: true,
		Errors:          validationErrors,
	}, template.lastVars)
}

func TestPostEditTeamOtherError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockEditTeamClient{}
	client.team.data = sirius.Team{
		DisplayName: "Complaints team",
	}
	client.editTeam.err = expectedErr
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/edit/123", nil)

	err := editTeam(client, template)(client.requiredPermissions(), w, r)

	assert.Equal(expectedErr, err)

	assert.Equal(0, template.count)
}

func TestPostEditTeamRetrieveError(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditTeamClient{}
	client.team.err = StatusError(http.StatusNotFound)
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/edit/123", nil)

	err := editTeam(client, template)(client.requiredPermissions(), w, r)

	assert.Equal(StatusError(http.StatusNotFound), err)

	assert.Equal(0, client.editTeam.count)
	assert.Equal(0, template.count)
}

func TestPostEditTeamRefDataError(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditTeamClient{}
	client.teamTypes.err = StatusError(http.StatusNotFound)
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/edit/123", nil)

	err := editTeam(client, template)(client.requiredPermissions(), w, r)

	assert.Equal(StatusError(http.StatusNotFound), err)

	assert.Equal(0, client.editTeam.count)
	assert.Equal(0, template.count)
}

func TestBadMethodEditTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditTeamClient{}
	client.team.data = sirius.Team{
		DisplayName: "Complaints team",
	}
	template := &mockTemplate{}

	r, _ := http.NewRequest("DELETE", "/teams/edit/123", nil)

	err := editTeam(client, template)(client.requiredPermissions(), nil, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.editTeam.count)
	assert.Equal(0, template.count)
}
