package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type WorkflowInformation interface {
	GetCurrentUserDetails(sirius.Context) (sirius.UserDetails, error)
	GetTaskTypes(sirius.Context, []string) ([]sirius.ApiTaskTypes, error)
	GetTaskList(sirius.Context, int, int, int, int, []string, []sirius.ApiTaskTypes, []string) (sirius.TaskList, int, error)
	GetPageDetails(sirius.TaskList, int, int) sirius.PageDetails
	GetTeamsForSelection(sirius.Context, int, []string) ([]sirius.ReturnedTeamCollection, error)
	GetAssigneesForFilter(sirius.Context, int, []string) (sirius.AssigneesTeam, error)
	AssignTasksToCaseManager(sirius.Context, int, string) error
	GetAppliedFilters(int, []sirius.ApiTaskTypes, []sirius.ReturnedTeamCollection, sirius.AssigneesTeam) []string
}

type workflowVars struct {
	Path           string
	XSRFToken      string
	MyDetails      sirius.UserDetails
	TaskList       sirius.TaskList
	PageDetails    sirius.PageDetails
	LoadTasks      []sirius.ApiTaskTypes
	TeamSelection  []sirius.ReturnedTeamCollection
	Assignees      sirius.AssigneesTeam
	AppliedFilters []string
	SuccessMessage string
	Error          string
	Errors         sirius.ValidationErrors
}

func checkForChangesToSelectedPagination(bothDisplayTaskLimits []string, currentTaskDisplayString string) int {
	currentTaskDisplay, _ := strconv.Atoi(currentTaskDisplayString)

	if len(bothDisplayTaskLimits) != 0 {
		topDisplayTaskLimit, _ := strconv.Atoi(bothDisplayTaskLimits[0])
		bottomDisplayTaskLimit, _ := strconv.Atoi(bothDisplayTaskLimits[1])
		if topDisplayTaskLimit != currentTaskDisplay {
			return topDisplayTaskLimit
		} else if bottomDisplayTaskLimit != currentTaskDisplay {
			return bottomDisplayTaskLimit
		} else {
			return currentTaskDisplay
		}
	}
	return 25
}

func getLoggedInTeam(myDetails sirius.UserDetails, defaultWorkflowTeam int) int {
	if len(myDetails.Teams) < 1 {
		return defaultWorkflowTeam
	} else {
		return myDetails.Teams[0].TeamId
	}
}

func getAssigneeIdForTask(teamId, assigneeId string) (int, error) {
	var assigneeIdForTask int
	var err error

	if assigneeId != "" {
		assigneeIdForTask, err = strconv.Atoi(assigneeId)
	} else if teamId != "" {
		assigneeIdForTask, err = strconv.Atoi(teamId)
	}
	if err != nil {
		return 0, err
	}
	return assigneeIdForTask, nil
}

func createTaskIdForUrl(taskIdArray []string) string {
	taskIdForUrl := ""

	for i := 0; i < len(taskIdArray); i++ {
		taskIdForUrl += taskIdArray[i]
		if i < (len(taskIdArray) - 1) {
			taskIdForUrl += "+"
		}
	}
	return taskIdForUrl
}

func loggingInfoForWorkflow(client WorkflowInformation, tmpl Template, defaultWorkflowTeam int) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		search, _ := strconv.Atoi(r.FormValue("page"))
		selectedTeamId, _ := strconv.Atoi(r.FormValue("change-team"))

		displayTaskLimit := checkForChangesToSelectedPagination(r.Form["tasksPerPage"], r.FormValue("currentTaskDisplay"))

		err := r.ParseForm()
		if err != nil {
			return err
		}

		taskTypeSelected := r.Form["selected-task-type"]
		assigneeSelected := r.Form["selected-assignee"]

		myDetails, err := client.GetCurrentUserDetails(ctx)
		if err != nil {
			return err
		}

		loggedInTeamId := getLoggedInTeam(myDetails, defaultWorkflowTeam)

		loadTaskTypes, err := client.GetTaskTypes(ctx, taskTypeSelected)
		if err != nil {
			return err
		}

		taskList, teamId, err := client.GetTaskList(ctx, search, displayTaskLimit, selectedTeamId, loggedInTeamId, taskTypeSelected, loadTaskTypes, assigneeSelected)
		if err != nil {
			return err
		}

		pageDetails := client.GetPageDetails(taskList, search, displayTaskLimit)

		teamSelection, err := client.GetTeamsForSelection(ctx, teamId, assigneeSelected)
		if err != nil {
			return err
		}

		assigneesForFilter, err := client.GetAssigneesForFilter(ctx, teamId, assigneeSelected)
		if err != nil {
			return err
		}

		appliedFilters := client.GetAppliedFilters(teamId, loadTaskTypes, teamSelection, assigneesForFilter)

		vars := workflowVars{
			Path:           r.URL.Path,
			XSRFToken:      ctx.XSRFToken,
			MyDetails:      myDetails,
			TaskList:       taskList,
			PageDetails:    pageDetails,
			LoadTasks:      loadTaskTypes,
			TeamSelection:  teamSelection,
			Assignees:      assigneesForFilter,
			AppliedFilters: appliedFilters,
		}

		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		case http.MethodPost:
			var newAssigneeIdForTask int
			selectedTeamToAssignTaskString := r.FormValue("assignTeam")
			if selectedTeamToAssignTaskString == "0" {
				vars.Errors = sirius.ValidationErrors{
					"selection": {"": "Please select a team"},
				}

				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			//this is where it picks up the new user to assign task to
			newAssigneeIdForTask, err = getAssigneeIdForTask(selectedTeamToAssignTaskString, r.FormValue("assignCM"))
			if err != nil {
				return err
			}

			err := r.ParseForm()
			if err != nil {
				return err
			}

			taskIdArray := (r.Form["selected-tasks"])
			taskIdForUrl := createTaskIdForUrl(taskIdArray)

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

			vars.TaskList, _, err = client.GetTaskList(ctx, search, displayTaskLimit, selectedTeamId, loggedInTeamId, taskTypeSelected, loadTaskTypes, assigneeSelected)
			if err != nil {
				return err
			}

			vars.PageDetails = client.GetPageDetails(taskList, search, displayTaskLimit)

			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
