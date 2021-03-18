package server

import (
	"net/http"

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

func listTaskTypes(client GetTaskTypeClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		// if !perm.HasPermission("team", http.MethodPut) {
		// 	return StatusError(http.StatusForbidden)
		// }

		// if r.Method != http.MethodGet {
		// 	return StatusError(http.StatusMethodNotAllowed)
		// }

		ctx := getContext(r)

		loadTaskTypes, err := client.GetTaskDetails(ctx)
		if err != nil {
			return err
		}

		vars := listTeamsVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			LoadTaskTypes: loadTaskTypes,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
