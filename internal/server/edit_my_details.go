package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type EditMyDetailsClient interface {
	MyDetails(sirius.Context) (sirius.MyDetails, error)
	EditMyDetails(sirius.Context, int, string) error
}

type editMyDetailsVars struct {
	Path        string
	XSRFToken   string
	Success     bool
	Errors      sirius.ValidationErrors
	PhoneNumber string
}

func editMyDetails(client EditMyDetailsClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if !perm.HasPermission("v1-users-updatetelephonenumber", http.MethodPut) {
			return StatusError(http.StatusForbidden)
		}

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.MyDetails(ctx)
		if err != nil {
			return err
		}

		vars := editMyDetailsVars{
			Path:        r.URL.Path,
			XSRFToken:   ctx.XSRFToken,
			PhoneNumber: myDetails.PhoneNumber,
		}

		if r.Method == http.MethodPost {
			vars.PhoneNumber = r.FormValue("phonenumber")
			err := client.EditMyDetails(ctx, myDetails.ID, vars.PhoneNumber)

			if e, ok := err.(*sirius.ValidationError); ok {
				vars.Errors = e.Errors
				w.WriteHeader(http.StatusBadRequest)
			} else if err != nil {
				return err
			} else {
				vars.Success = true
			}
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
