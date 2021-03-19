package server

import (
	"net/http"
	"sort"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type GetTaskTypeClient interface {
	GetTaskDetails(sirius.Context) ([]sirius.ApiTaskTypes, error)
}

type listTaskTypeVars struct {
	Path      string
	XSRFToken string
	LoadTasks []sirius.ApiTaskTypes
}

func listTaskTypes(client GetTaskTypeClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		loadTaskTypes, err := client.GetTaskDetails(ctx)
		if err != nil {
			return err
		}

		sort.SliceStable(loadTaskTypes, func(i, j int) bool {
			return loadTaskTypes[i].Incomplete < loadTaskTypes[j].Incomplete
		})

		vars := listTaskTypeVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
			LoadTasks: loadTaskTypes,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
