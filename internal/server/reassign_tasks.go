package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type ReassignTasksClient interface {
	ReassignTasks(sirius.Context, sirius.ReassignTasksParams) (string, error)
}

func reassignTasks(client ReassignTasksClient) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

		err := r.ParseForm()
		if err != nil {
			return err
		}

		reassignSuccessMessage, err := client.ReassignTasks(ctx, sirius.ReassignTasksParams{
			AssignTeam: r.FormValue("assignTeam"),
			AssignCM:   r.FormValue("assignCM"),
			TaskIds:    r.Form["selected-tasks"],
			IsPriority: r.FormValue("priority"),
		})

		if err != nil {
			return err
		}

		return Redirect{
			Path:           r.URL.RequestURI(),
			SuccessMessage: reassignSuccessMessage,
		}
	}
}
