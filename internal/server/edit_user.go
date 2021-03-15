package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type EditUserClient interface {
	User(sirius.Context, int) (sirius.AuthUser, error)
	EditUser(sirius.Context, sirius.AuthUser) error
	Roles(sirius.Context) ([]string, error)
}

type editUserVars struct {
	Path              string
	XSRFToken         string
	Roles             []string
	User              sirius.AuthUser
	Success           bool
	Errors            sirius.ValidationErrors
}

func editUser(client EditUserClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if !perm.HasPermission("v1-users", http.MethodPut) {
			return StatusError(http.StatusForbidden)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/edit-user/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		ctx := getContext(r)

		roles, err := client.Roles(ctx)
		if err != nil {
			return err
		}

		vars := editUserVars{
			Path:              r.URL.Path,
			XSRFToken:         ctx.XSRFToken,
			Roles:             roles,
		}

		switch r.Method {
		case http.MethodGet:
			user, err := client.User(ctx, id)
			if err != nil {
				return err
			}
			vars.User = user

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			vars.User = sirius.AuthUser{
				ID:           id,
				Firstname:    r.PostFormValue("firstname"),
				Surname:      r.PostFormValue("surname"),
				Organisation: r.PostFormValue("organisation"),
				Roles:        r.PostForm["roles"],
				Suspended:    r.PostFormValue("suspended") == "Yes",
				Locked:       r.PostFormValue("locked") == "Yes",
			}
			err := client.EditUser(ctx, vars.User)

			if _, ok := err.(sirius.ClientError); ok {
				vars.Errors = sirius.ValidationErrors{
					"firstname": {
						"": err.Error(),
					},
				}

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if err != nil {
				return err
			}

			vars.Success = true
			vars.User.Email = r.PostFormValue("email")
			return tmpl.ExecuteTemplate(w, "page", vars)

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
