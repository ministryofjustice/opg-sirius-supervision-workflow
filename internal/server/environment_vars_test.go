package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvironmentVars(t *testing.T) {
	vars, err := NewEnvironmentVars()

	assert.Nil(t, err)
	assert.Equal(t, EnvironmentVars{
		Port:                  "1234",
		WebDir:                "web",
		SiriusURL:             "http://localhost:8080",
		SiriusPublicURL:       "",
		Prefix:                "/supervision/workflow",
		DefaultWorkflowTeamID: 21,
		DefaultPaTeamID:       "28",
		DefaultProTeamID:      "31",
	}, vars)
}
