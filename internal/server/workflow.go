package server

import (
	"errors"
	"fmt"
	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"net/http"
	"os"
	"strconv"
)

type WorkflowInformation interface {
	GetCurrentUserDetails(sirius.Context) (sirius.UserDetails, error)
	GetTaskTypes(sirius.Context, []string) ([]sirius.ApiTaskTypes, error)
	GetTaskList(sirius.Context, int, int, sirius.ReturnedTeamCollection, []string, []sirius.ApiTaskTypes, []string) (sirius.TaskList, error)
	GetPageDetails(sirius.TaskList, int, int) sirius.PageDetails
	GetTeamsForSelection(sirius.Context) ([]sirius.ReturnedTeamCollection, error)
	AssignTasksToCaseManager(sirius.Context, int, string) error
}

type workflowVars struct {
	Path               string
	XSRFToken          string
	MyDetails          sirius.UserDetails
	TaskList           sirius.TaskList
	PageDetails        sirius.PageDetails
	LoadTasks          []sirius.ApiTaskTypes
	TeamSelection      []sirius.ReturnedTeamCollection
	SelectedTeam       sirius.ReturnedTeamCollection
	SelectedAssignees  []string
	SelectedUnassigned string
	AppliedFilters     []string
	SuccessMessage     string
	Error              string
	Errors             sirius.ValidationErrors
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

func getLoggedInTeamId(myDetails sirius.UserDetails, defaultWorkflowTeam int) int {
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
		logger.Print("getAssigneeIdForTask error: " + err.Error())
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

func getSelectedTeam(r *http.Request, loggedInTeamId int, defaultTeamId int, teamSelection []sirius.ReturnedTeamCollection) (sirius.ReturnedTeamCollection, error) {
	selectors := []string{
		r.URL.Query().Get("change-team"),
		r.FormValue("change-team"),
		strconv.Itoa(loggedInTeamId),
		strconv.Itoa(defaultTeamId),
	}

	for _, selector := range selectors {
		for _, team := range teamSelection {
			if team.Selector == selector {
				return team, nil
			}
		}
	}

	return sirius.ReturnedTeamCollection{}, errors.New("invalid team selection")
}

func loggingInfoForWorkflow(client WorkflowInformation, tmpl Template, defaultWorkflowTeam int) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logging.New(os.Stdout, "opg-sirius-workflow ")
		ctx := getContext(r)
		search, _ := strconv.Atoi(r.FormValue("page"))
		if search < 1 {
			search = 1
		}

		err := r.ParseForm()
		if err != nil {
			logger.Print("ParseForm error: " + err.Error())
			return err
		}

		displayTaskLimit := checkForChangesToSelectedPagination(r.Form["tasksPerPage"], r.FormValue("currentTaskDisplay"))

		myDetails, err := client.GetCurrentUserDetails(ctx)
		if err != nil {
			logger.Print("GetCurrentUserDetails error " + err.Error())
			return err
		}

		teamSelection, err := client.GetTeamsForSelection(ctx)
		if err != nil {
			logger.Print("GetTeamsForSelection error " + err.Error())
			return err
		}

		loggedInTeamId := getLoggedInTeamId(myDetails, defaultWorkflowTeam)

		selectedTeam, err := getSelectedTeam(r, loggedInTeamId, defaultWorkflowTeam, teamSelection)
		if err != nil {
			logger.Print("getSelectedTeam error " + err.Error())
			return err
		}

		selectedAssignees := r.Form["selected-assignee"]
		selectedUnassigned := r.FormValue("selected-unassigned")

		if selectedUnassigned == selectedTeam.Selector {
			selectedAssignees = append(selectedAssignees, strconv.Itoa(selectedTeam.Id))
			for _, t := range selectedTeam.Teams {
				selectedAssignees = append(selectedAssignees, strconv.Itoa(t.Id))
			}
		}

		selectedTaskType := r.Form["selected-task-type"]

		taskTypes, err := client.GetTaskTypes(ctx, selectedTaskType)
		if err != nil {
			logger.Print("GetTaskTypes error " + err.Error())
			return err
		}

		taskList, err := client.GetTaskList(ctx, search, displayTaskLimit, selectedTeam, selectedTaskType, taskTypes, selectedAssignees)

		if err != nil {
			logger.Print("GetTaskList error " + err.Error())
			return err
		}
		if search > taskList.Pages.PageTotal && taskList.Pages.PageTotal > 0 {
			search = taskList.Pages.PageTotal
			taskList, err = client.GetTaskList(ctx, search, displayTaskLimit, selectedTeam, selectedTaskType, taskTypes, selectedAssignees)

			if err != nil {
				logger.Print("GetTaskList error " + err.Error())
				return err
			}
		}

		pageDetails := client.GetPageDetails(taskList, search, displayTaskLimit)

		appliedFilters := sirius.GetAppliedFilters(selectedTeam, selectedAssignees, selectedUnassigned, taskTypes)

		vars := workflowVars{
			Path:               r.URL.Path,
			XSRFToken:          ctx.XSRFToken,
			MyDetails:          myDetails,
			TaskList:           taskList,
			PageDetails:        pageDetails,
			LoadTasks:          taskTypes,
			TeamSelection:      teamSelection,
			SelectedTeam:       selectedTeam,
			SelectedAssignees:  selectedAssignees,
			SelectedUnassigned: selectedUnassigned,
			AppliedFilters:     appliedFilters,
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
					"selection": map[string]string{"": "Please select a team"},
				}

				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			//this is where it picks up the new user to assign task to
			newAssigneeIdForTask, err = getAssigneeIdForTask(logger, selectedTeamToAssignTaskString, r.FormValue("assignCM"))
			if err != nil {
				logger.Print("getAssigneeIdForTask error: " + err.Error())
				return err
			}

			err := r.ParseForm()
			if err != nil {
				logger.Print("ParseForm error: " + err.Error())
				return err
			}

			taskIdArray := r.Form["selected-tasks"]
			taskIdForUrl := createTaskIdForUrl(taskIdArray)

			if err != nil {
				logger.Print("taskIdForUrl error: " + err.Error())
				return err
			}

			// Attempt to save
			err = client.AssignTasksToCaseManager(ctx, newAssigneeIdForTask, taskIdForUrl)
			if err != nil {
				logger.Print("AssignTasksToCaseManager: " + err.Error())
				return err
			}

			if vars.Errors == nil {
				vars.SuccessMessage = fmt.Sprintf("%d tasks have been reassigned", len(taskIdArray))
			}

			vars.TaskList, err = client.GetTaskList(ctx, search, displayTaskLimit, selectedTeam, selectedTaskType, taskTypes, selectedAssignees)
			if err != nil {
				logger.Print("vars.TaskList error: " + err.Error())
				return err
			}

			vars.PageDetails = client.GetPageDetails(vars.TaskList, search, displayTaskLimit)

			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
