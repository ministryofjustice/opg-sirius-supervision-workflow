package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type WorkflowInformation interface {
	SiriusUserDetails(sirius.Context) (sirius.UserDetails, error)
	GetTaskType(sirius.Context, []string) ([]sirius.ApiTaskTypes, error)
	GetTaskList(sirius.Context, int, int, int, int, []string, []sirius.ApiTaskTypes, []string) (sirius.TaskList, error)
	GetTaskDetails(sirius.Context, sirius.TaskList, int, int) sirius.TaskDetails
	GetTeamSelection(sirius.Context, int, int) ([]sirius.ReturnedTeamCollection, error)
	GetAssigneesForFilter(sirius.Context, int, int, []string) (sirius.AssigneesTeam, error)
	AssignTasksToCaseManager(sirius.Context, int, string) error
}

type workflowVars struct {
	Path           string
	XSRFToken      string
	MyDetails      sirius.UserDetails
	TaskList       sirius.TaskList
	TaskDetails    sirius.TaskDetails
	LoadTasks      []sirius.ApiTaskTypes
	TeamSelection  []sirius.ReturnedTeamCollection
	Assignees      sirius.AssigneesTeam
	SuccessMessage string
	Error          string
	Errors         sirius.ValidationErrors
}

func loggingInfoForWorflow(client WorkflowInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		var displayTaskLimit int

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		search, _ := strconv.Atoi(r.FormValue("page"))
		bothDisplayTaskLimits := r.Form["tasksPerPage"]
		currentTaskDisplay, _ := strconv.Atoi(r.FormValue("currentTaskDisplay"))

		if len(bothDisplayTaskLimits) != 0 {
			topDisplayTaskLimit, _ := strconv.Atoi(bothDisplayTaskLimits[0])
			bottomDisplayTaskLimit, _ := strconv.Atoi(bothDisplayTaskLimits[1])
			if topDisplayTaskLimit != currentTaskDisplay {
				displayTaskLimit = topDisplayTaskLimit
			} else if bottomDisplayTaskLimit != currentTaskDisplay {
				displayTaskLimit = bottomDisplayTaskLimit
			} else {
				displayTaskLimit = currentTaskDisplay
			}
		} else {
			displayTaskLimit = 25
		}

		selectedTeamName, _ := strconv.Atoi(r.FormValue("change-team"))
		selectedTeamToAssignTaskString := r.FormValue("assignTeam")

		err := r.ParseForm()
		if err != nil {
			return err
		}
		taskTypeSelected := (r.Form["selected-task-type"])
		assigneeSelected := (r.Form["selected-assignee"])

		myDetails, err := client.SiriusUserDetails(ctx)
		if err != nil {
			return err
		}

		if len(myDetails.Teams) < 1 {
			err = errors.New("This user is not in a team")
		}
		if err != nil {
			return err
		}

		loggedInTeamId := myDetails.Teams[0].TeamId
		loadTaskTypes, err := client.GetTaskType(ctx, taskTypeSelected)
		if err != nil {
			return err
		}

		taskList, err := client.GetTaskList(ctx, search, displayTaskLimit, selectedTeamName, loggedInTeamId, taskTypeSelected, loadTaskTypes, assigneeSelected)
		if err != nil {
			return err
		}

		taskdetails := client.GetTaskDetails(ctx, taskList, search, displayTaskLimit)

		teamSelection, err := client.GetTeamSelection(ctx, loggedInTeamId, selectedTeamName)
		if err != nil {
			return err
		}

		assigneesForFilter, err := client.GetAssigneesForFilter(ctx, teamSelection[0].Id, selectedTeamName, assigneeSelected)
		if err != nil {
			return err
		}

		vars := workflowVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			MyDetails:     myDetails,
			TaskList:      taskList,
			TaskDetails:   taskdetails,
			LoadTasks:     loadTaskTypes,
			TeamSelection: teamSelection,
			Assignees:     assigneesForFilter,
		}

		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		case http.MethodPost:

			if selectedTeamToAssignTaskString == "0" {
				vars.Errors = sirius.ValidationErrors{
					"selection": {"": "Please select a team"},
				}

				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			checkTaskHasIdForAssigning := r.PostFormValue("assignCM")
			var newAssigneeIdForTask int

			if checkTaskHasIdForAssigning != "" {
				newAssigneeIdForTask, _ = strconv.Atoi(r.PostFormValue("assignCM"))
			} else {
				newAssigneeIdForTask = vars.Assignees.Members[0].TeamMembersId
			}

			err := r.ParseForm()
			if err != nil {
				return err
			}
			taskIdArray := (r.Form["selected-tasks"])

			taskIdForUrl := ""

			for i := 0; i < len(taskIdArray); i++ {
				taskIdForUrl += taskIdArray[i]
				if i < (len(taskIdArray) - 1) {
					taskIdForUrl += "+"
				}
			}

			if err != nil {
				return err
			}

			// Attempt to save
			err = client.AssignTasksToCaseManager(ctx, newAssigneeIdForTask, taskIdForUrl)
			if err != nil {
				return err
			}
			if vars.Errors == nil {
				vars.SuccessMessage = fmt.Sprintf("%d tasks have been reassigned", len(taskIdArray))
			}
			TaskList, err := client.GetTaskList(ctx, search, displayTaskLimit, selectedTeamName, loggedInTeamId, taskTypeSelected, loadTaskTypes, assigneeSelected)
			if err != nil {
				return err
			}

			vars.TaskList = TaskList
			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
