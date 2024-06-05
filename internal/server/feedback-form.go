package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"net/http"
)

type FeedbackFormClient interface {
	GetCurrentUserDetails(sirius.Context) (model.Assignee, error)
}

func feedbackForm(client FeedbackFormClient, tmpl Template) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		return tmpl.Execute(w, app)
	}
}
