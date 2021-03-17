package server

import (
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type GetTaskTypeClient interface {
	GetTaskDetails(sirius.Context) ([]sirius.LoadTasks, error)
}

type listTeamsVars struct {
	Path          string
	XSRFToken     string
	LoadTaskTypes []sirius.LoadTasks
}

func listTaskTypes(client GetTaskTypeClient, tmpl Templates) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if !perm.HasPermission("team", http.MethodPut) {
			return StatusError(http.StatusForbidden)
		}

		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		teams, err := client.Teams(ctx)
		if err != nil {
			return err
		}

		search := r.FormValue("search")
		if search != "" {
			searchLower := strings.ToLower(search)

			var matchingTeams []sirius.Team
			for _, t := range teams {
				if strings.Contains(strings.ToLower(t.DisplayName), searchLower) {
					matchingTeams = append(matchingTeams, t)
				}
			}

			teams = matchingTeams
		}

		vars := listTeamsVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
			Search:    search,
			Teams:     teams,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
