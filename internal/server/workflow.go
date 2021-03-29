package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type UserDetailsClient interface {
	SiriusUserDetails(sirius.Context) (sirius.UserDetails, error)
	GetTaskDetails(sirius.Context) ([]sirius.ApiTaskTypes, error)
	GetTaskList(sirius.Context) (sirius.ApiTask, error)
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
	TaskList           sirius.ApiTask
	LoadTasks          []sirius.ApiTaskTypes
}

func loggingInfoForWorflow(client UserDetailsClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.SiriusUserDetails(ctx)
		loadTaskTypes, err := client.GetTaskDetails(ctx)
		taskList, err := client.GetTaskList(ctx)
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
			TaskList:    taskList,
			LoadTasks:   loadTaskTypes,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
