package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type ResendConfirmationClient interface {
	ResendConfirmation(sirius.Context, string) error
}

type resendConfirmationVars struct {
	Path  string
	ID    string
	Email string
}

func resendConfirmation(client ResendConfirmationClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if !perm.HasPermission("v1-users", http.MethodPut) {
			return StatusError(http.StatusForbidden)
		}

		switch r.Method {
		case http.MethodGet:
			return RedirectError("/users")

		case http.MethodPost:
			vars := resendConfirmationVars{
				Path:  r.URL.Path,
				ID:    r.PostFormValue("id"),
				Email: r.PostFormValue("email"),
			}

			err := client.ResendConfirmation(getContext(r), vars.Email)
			if err != nil {
				return err
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
