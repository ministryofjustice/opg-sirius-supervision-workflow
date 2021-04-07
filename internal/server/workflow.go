package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type WorkflowInformation interface {
	SiriusUserDetails(sirius.Context) (sirius.UserDetails, error)
	GetTaskType(sirius.Context) (sirius.TaskTypes, error)
	GetTaskList(sirius.Context) (sirius.TaskList, error)
}

type workflowVars struct {
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
	TaskList           sirius.TaskList
	LoadTasks          sirius.TaskTypes
}

func loggingInfoForWorflow(client WorkflowInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		myDetails, err := client.SiriusUserDetails(ctx)
		loadTaskTypes, err := client.GetTaskType(ctx)
		taskList, err := client.GetTaskList(ctx)
		if err != nil {
			return err
		}

		vars := workflowVars{
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
