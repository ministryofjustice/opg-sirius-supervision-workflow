package server

import (
	"errors"
	"os"
	"strconv"
)

type EnvironmentVars struct {
	Port                  string
	WebDir                string
	SiriusURL             string
	SiriusPublicURL       string
	Prefix                string
	DefaultWorkflowTeamID int
	DefaultPaTeamID       string
	DefaultProTeamID      string
	FinanceAdminLink      string
}

func NewEnvironmentVars() (EnvironmentVars, error) {
	defaultTeamId, err := strconv.Atoi(getEnv("DEFAULT_WORKFLOW_TEAM", "21"))
	if err != nil {
		return EnvironmentVars{}, errors.New("error converting DEFAULT_WORKFLOW_TEAM to int")
	}

	return EnvironmentVars{
		Port:                  getEnv("PORT", "1234"),
		WebDir:                getEnv("WEB_DIR", "web"),
		SiriusURL:             getEnv("SIRIUS_URL", "http://localhost:8080"),
		SiriusPublicURL:       getEnv("SIRIUS_PUBLIC_URL", ""),
		Prefix:                getEnv("PREFIX", ""),
		DefaultWorkflowTeamID: defaultTeamId,
		DefaultPaTeamID:       getEnv("DEFAULT_PA_ECM_TEAM", "28"),
		DefaultProTeamID:      getEnv("DEFAULT_PRO_ECM_TEAM", "31"),
		FinanceAdminLink:      getEnv("FINANCE_ADMIN_LINK", "0"),
	}, nil
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
