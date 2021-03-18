package server

import (
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type UserDetailsClient interface {
	SiriusUserDetails(sirius.Context) (sirius.UserDetails, error)
}

type userDetailsVars struct {
	Path               string
	ID                 int
	Firstname          string
	Surname            string
	Email              string
	PhoneNumber        string
	Organisation       string
	Roles              []string
	Teams              []string
	CanEditPhoneNumber bool
}

func loggingInfoForWorflow(client UserDetailsClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		log.Println("my details page response is")
		log.Println("request is")
		log.Println(r)
		log.Println("writer is")
		log.Println(w)
		log.Println("permissions is")
		log.Println(perm)

		ctx := getContext(r)

		myDetails, err := client.SiriusUserDetails(ctx)
		if err != nil {
			return err
		}

		vars := userDetailsVars{
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

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
