package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type EditTeamClient interface {
	Team(sirius.Context, int) (sirius.Team, error)
	EditTeam(sirius.Context, sirius.Team) error
	TeamTypes(sirius.Context) ([]sirius.RefDataTeamType, error)
}

type editTeamVars struct {
	Path            string
	XSRFToken       string
	Team            sirius.Team
	TeamTypeOptions []sirius.RefDataTeamType
	CanEditTeamType bool
	CanDeleteTeam   bool
	Success         bool
	Errors          sirius.ValidationErrors
}

func editTeam(client EditTeamClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if !perm.HasPermission("team", http.MethodPut) {
			return StatusError(http.StatusForbidden)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/edit/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		ctx := getContext(r)

		team, err := client.Team(ctx, id)
		if err != nil {
			return err
		}

		canEditTeamType := perm.HasPermission("team", http.MethodPost)
		canDeleteTeam := perm.HasPermission("v1-teams", http.MethodDelete)

		teamTypes, err := client.TeamTypes(ctx)
		if err != nil {
			return err
		}

		vars := editTeamVars{
			Path:            r.URL.Path,
			XSRFToken:       ctx.XSRFToken,
			Team:            team,
			TeamTypeOptions: teamTypes,
			CanEditTeamType: canEditTeamType,
			CanDeleteTeam:   canDeleteTeam,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		case http.MethodPost:
			vars.Team.DisplayName = r.PostFormValue("name")
			vars.Team.PhoneNumber = r.PostFormValue("phone")
			vars.Team.Email = r.PostFormValue("email")

			if canEditTeamType {
				if r.PostFormValue("service") == "supervision" {
					vars.Team.Type = r.PostFormValue("supervision-type")
				} else {
					vars.Team.Type = ""
				}
			} else {
				vars.Team.Type = team.Type
			}

			// Attempt to save
			err := client.EditTeam(ctx, vars.Team)

			if e, ok := err.(*sirius.ValidationError); ok {
				vars.Errors = e.Errors
				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				return err
			} else {
				vars.Success = true
			}

			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
