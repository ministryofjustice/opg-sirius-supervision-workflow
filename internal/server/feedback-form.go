package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"net/http"
)

type FeedbackFormClient interface {
	SubmitFeedback(sirius.Context, model.FeedbackForm) error
}

func feedbackForm(client FeedbackFormClient, tmpl Template) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}
		ctx := getContext(r)

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return err
			}

			err = client.SubmitFeedback(ctx, model.FeedbackForm{
				Id:         app.MyDetails.Id,
				Name:       r.FormValue("name"),
				Email:      r.FormValue("email"),
				CaseNumber: r.FormValue("case-number"),
				Feedback:   r.FormValue("more-detail"),
			})

			//if err != nil {
			//	return err
			//}

			app.SuccessMessage = "Form Submitted"
		}

		return tmpl.Execute(w, app)
	}
}
