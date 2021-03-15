package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type DeleteTeamClient interface {
	Team(sirius.Context, int) (sirius.Team, error)
	DeleteTeam(sirius.Context, int) error
}

type deleteTeamVars struct {
	Path           string
	XSRFToken      string
	Team           sirius.Team
	Errors         sirius.ValidationErrors
	SuccessMessage string
}

func deleteTeam(client DeleteTeamClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if !perm.HasPermission("v1-teams", http.MethodDelete) {
			return StatusError(http.StatusForbidden)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/delete/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		team, err := client.Team(ctx, id)
		if err != nil {
			return err
		}

		vars := deleteTeamVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
			Team:      team,
		}

		if r.Method == http.MethodPost {
			err := client.DeleteTeam(ctx, id)

			if _, ok := err.(sirius.ClientError); ok {
				vars.Errors = sirius.ValidationErrors{
					"": {
						"": err.Error(),
					},
				}

				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				return err
			} else {
				vars.SuccessMessage = fmt.Sprintf("The team \"%s\" was deleted.", team.DisplayName)
			}
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
