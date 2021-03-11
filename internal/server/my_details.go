package server

import (
	"context"
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type myDetailsClient interface {
	MyDetails(context.Context, []*http.Cookie) (sirius.MyDetails, error)
	AuthenticateClient
}

type myDetailsVars struct {
	Path         string
	ID           int
	Firstname    string
	Surname      string
	Email        string
	PhoneNumber  string
	Organisation string
	Roles        []string
	Teams        []string
}

func myDetails(logger *log.Logger, client MyDetailsClient, templates Templates) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		myDetails, err := client.MyDetails(r.Context(), r.Cookies())
		if err == sirius.ErrUnauthorized {
			client.Authenticate(w, r)
			return
		} else if err != nil {
			logger.Println("myDetails:", err)
			http.Error(w, "Could not connect to Sirius", http.StatusInternalServerError)
			return
		}

		vars := myDetailsVars{
			Path:        r.URL.Path,
			ID:          myDetails.ID,
			Firstname:   myDetails.Firstname,
			Surname:     myDetails.Surname,
			Email:       myDetails.Email,
			PhoneNumber: myDetails.PhoneNumber,
		}

		for _, role := range myDetails.Roles {
			if role == "OPG User" || role == "COP User" {
				vars.Organisation = role
			} else {
				vars.Roles = append(vars.Roles, role)
			}
		}

		for _, team := range myDetails.Teams {
			vars.Teams = append(vars.Teams, team.DisplayName)
		}

		if err := templates.ExecuteTemplate(w, "my-details.gotmpl", vars); err != nil {
			logger.Println("myDetails:", err)
		}
	})
}
