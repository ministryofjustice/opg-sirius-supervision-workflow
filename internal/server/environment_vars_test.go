package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEnvironmentVars(t *testing.T) {
	vars, err := NewEnvironmentVars()

	assert.Nil(t, err)
	assert.Equal(t, EnvironmentVars{
		Port:            "1234",
		WebDir:          "web",
		SiriusURL:       "http://localhost:8080",
		SiriusPublicURL: "",
		Prefix:          "",
		DefaultTeamId:   21,
		ShowDeputies:    false,
	}, vars)
}
