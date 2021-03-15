package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

type UnlockUserClient interface {
	User(sirius.Context, int) (sirius.AuthUser, error)
	EditUser(sirius.Context, sirius.AuthUser) error
}

type unlockUserVars struct {
	Path      string
	XSRFToken string
	User      sirius.AuthUser
	Errors    sirius.ValidationErrors
}

func unlockUser(client UnlockUserClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if !perm.HasPermission("v1-users", http.MethodPut) {
			return StatusError(http.StatusForbidden)
		}

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/unlock-user/"))
		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		user, err := client.User(ctx, id)
		if err != nil {
			return err
		}

		vars := unlockUserVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
			User:      user,
		}

		if r.Method == http.MethodPost {
			user.Locked = false
			err := client.EditUser(ctx, user)

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
				return RedirectError(fmt.Sprintf("/edit-user/%d", user.ID))
			}
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
