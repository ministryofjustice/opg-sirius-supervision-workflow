package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/logging"
	"net/http"
	"os"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type WorkflowInformation interface {
	GetCurrentUserDetails(sirius.Context, *logging.Logger) (sirius.UserDetails, error)
	GetTaskTypes(sirius.Context, *logging.Logger, []string) ([]sirius.ApiTaskTypes, error)
	GetTaskList(sirius.Context, *logging.Logger, int, int, int, int, []string, []sirius.ApiTaskTypes, []string) (sirius.TaskList, int, error)
	GetPageDetails(*logging.Logger, sirius.TaskList, int, int) sirius.PageDetails
	GetTeamsForSelection(sirius.Context, *logging.Logger, int, []string) ([]sirius.ReturnedTeamCollection, error)
	GetAssigneesForFilter(sirius.Context, *logging.Logger, int, []string) (sirius.AssigneesTeam, error)
	AssignTasksToCaseManager(sirius.Context, *logging.Logger, int, string) error
	GetAppliedFilters(*logging.Logger, int, []sirius.ApiTaskTypes, []sirius.ReturnedTeamCollection, sirius.AssigneesTeam) []string
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

func getAssigneeIdForTask(logger *logging.Logger, teamId, assigneeId string) (int, error) {
	var assigneeIdForTask int
	var err error

	if assigneeId != "" {
		assigneeIdForTask, err = strconv.Atoi(assigneeId)
	} else if teamId != "" {
		assigneeIdForTask, err = strconv.Atoi(teamId)
	}
	if err != nil {
		logger.Print("getAssigneeIdForTask error")
		logger.Print(err)
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
		logger := logging.New(os.Stdout, "opg-sirius-workflow ")
		ctx := getContext(r)
		search, _ := strconv.Atoi(r.FormValue("page"))
		if search < 1 {
			search = 1
		}
		selectedTeamId, _ := strconv.Atoi(r.FormValue("change-team"))
		displayTaskLimit := checkForChangesToSelectedPagination(r.Form["tasksPerPage"], r.FormValue("currentTaskDisplay"))

		err := r.ParseForm()
		if err != nil {
			logger.Print("ParseForm error")
			logger.Print(err)
			return err
		}

		taskTypeSelected := r.Form["selected-task-type"]
		assigneeSelected := r.Form["selected-assignee"]

		myDetails, err := client.GetCurrentUserDetails(ctx, logger)
		if err != nil {
			logger.Print("GetCurrentUserDetails error")
			logger.Print(err)
			return err
		}

		loggedInTeamId := getLoggedInTeam(myDetails, defaultWorkflowTeam)

		loadTaskTypes, err := client.GetTaskTypes(ctx, logger, taskTypeSelected)
		if err != nil {
			logger.Print("GetTaskTypes error")
			logger.Print(err)
			return err
		}

		taskList, teamId, err := client.GetTaskList(ctx, logger, search, displayTaskLimit, selectedTeamId, loggedInTeamId, taskTypeSelected, loadTaskTypes, assigneeSelected)
		if err != nil {
			logger.Print("GetTaskList error")
			logger.Print(err)
			return err
		}

		pageDetails := client.GetPageDetails(logger, taskList, search, displayTaskLimit)

		teamSelection, err := client.GetTeamsForSelection(ctx, logger, teamId, assigneeSelected)
		if err != nil {
			logger.Print("GetTeamsForSelection error")
			logger.Print(err)
			return err
		}

		assigneesForFilter, err := client.GetAssigneesForFilter(ctx, logger, teamId, assigneeSelected)
		if err != nil {
			logger.Print("GetAssigneesForFilter error")
			logger.Print(err)
			return err
		}

		appliedFilters := client.GetAppliedFilters(logger, teamId, loadTaskTypes, teamSelection, assigneesForFilter)

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
			newAssigneeIdForTask, err = getAssigneeIdForTask(logger, selectedTeamToAssignTaskString, r.FormValue("assignCM"))
			if err != nil {
				logger.Print("getAssigneeIdForTask error")
				logger.Print(err)
				return err
			}

			err := r.ParseForm()
			if err != nil {
				logger.Print("ParseForm error")
				logger.Print(err)
				return err
			}

			taskIdArray := (r.Form["selected-tasks"])
			taskIdForUrl := createTaskIdForUrl(taskIdArray)

			if err != nil {
				logger.Print("taskIdForUrl error")
				logger.Print(err)
				return err
			}

			// Attempt to save
			err = client.AssignTasksToCaseManager(ctx, logger, newAssigneeIdForTask, taskIdForUrl)
			if err != nil {
				logger.Print("AssignTasksToCaseManager error")
				logger.Print(err)
				return err
			}

			if vars.Errors == nil {
				vars.SuccessMessage = fmt.Sprintf("%d tasks have been reassigned", len(taskIdArray))
			}

			vars.TaskList, _, err = client.GetTaskList(ctx, logger, search, displayTaskLimit, selectedTeamId, loggedInTeamId, taskTypeSelected, loadTaskTypes, assigneeSelected)
			if err != nil {
				logger.Print("vars.TaskList error")
				logger.Print(err)
				return err
			}

			vars.PageDetails = client.GetPageDetails(logger, vars.TaskList, search, displayTaskLimit)

			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
