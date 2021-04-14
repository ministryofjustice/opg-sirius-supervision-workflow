package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type WorkflowInformation interface {
	SiriusUserDetails(sirius.Context) (sirius.UserDetails, error)
	GetTaskType(sirius.Context) (sirius.TaskTypes, error)
	GetTaskList(sirius.Context, int, int) (sirius.TaskList, sirius.TaskDetails, error)
	GetTeamSelection(sirius.Context) ([]sirius.TeamCollection, error)
	GetTeamSelected(sirius.Context, int) (sirius.TeamSelected, error)
}

type workflowVars struct {
	Path          string
	MyDetails     sirius.UserDetails
	TaskList      sirius.TaskList
	TaskDetails   sirius.TaskDetails
	LoadTasks     sirius.TaskTypes
	TeamSelection []sirius.TeamCollection
	TeamSelected  sirius.TeamSelected
}

func loggingInfoForWorflow(client WorkflowInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		search, _ := strconv.Atoi(r.FormValue("page"))
		displayTaskLimit, _ := strconv.Atoi(r.FormValue("tasksPerPage"))
		selectedTeamName, _ := strconv.Atoi(r.FormValue("change-team"))

		myDetails, err := client.SiriusUserDetails(ctx)
		loadTaskTypes, err := client.GetTaskType(ctx)
		taskList, taskdetails, err := client.GetTaskList(ctx, search, displayTaskLimit)
		teamSelection, err := client.GetTeamSelection(ctx)
		teamSelected, err := client.GetTeamSelected(ctx, selectedTeamName)
		if err != nil {
			return err
		}

		vars := workflowVars{
			Path:          r.URL.Path,
			MyDetails:     myDetails,
			TaskList:      taskList,
			TaskDetails:   taskdetails,
			LoadTasks:     loadTaskTypes,
			TeamSelection: teamSelection,
			TeamSelected:  teamSelected,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
