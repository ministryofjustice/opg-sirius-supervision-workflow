package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type WorkflowInformation interface {
	SiriusUserDetails(sirius.Context) (sirius.UserDetails, error)
	GetTaskType(sirius.Context) (sirius.TaskTypes, error)
	GetTaskList(sirius.Context, int, int, sirius.TeamSelected) (sirius.TaskList, sirius.TaskDetails, error)
	GetTeamSelection(sirius.Context, sirius.UserDetails, int, int) ([]sirius.TeamCollection, error)
	GetTeamSelected(sirius.Context, []sirius.TeamCollection) (sirius.TeamSelected, error)
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
		oldTeamId, _ := strconv.Atoi(r.FormValue("feed-in-old-team-id"))
		// get workflow to submit every time page reloads
		log.Print("workflow selected team name")
		log.Print(selectedTeamName)

		myDetails, err := client.SiriusUserDetails(ctx)
		teamSelection, err := client.GetTeamSelection(ctx, myDetails, selectedTeamName, oldTeamId)
		selectedTeamMembers, err := client.GetTeamSelected(ctx, teamSelection)

		loadTaskTypes, err := client.GetTaskType(ctx)
		taskList, taskdetails, err := client.GetTaskList(ctx, search, displayTaskLimit, selectedTeamMembers)

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
			TeamSelected:  selectedTeamMembers,
		}

		// log.Print(vars.TeamSelected)

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
