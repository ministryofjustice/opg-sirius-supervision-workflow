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
	AssignTasksToCaseManager(sirius.Context, int, []string) error
}

func getTasksPerPage(valueFromUrl string) int {
	validOptions := []int{25, 50, 100}
	tasksPerPage, _ := strconv.Atoi(valueFromUrl)
	for _, opt := range validOptions {
		if opt == tasksPerPage {
			return tasksPerPage
		}
	}
	return validOptions[0]
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

func getSelectedTeam(r *http.Request, loggedInTeamId int, defaultTeamId int, teamSelection []sirius.ReturnedTeamCollection) (sirius.ReturnedTeamCollection, error) {
	selectors := []string{
		r.URL.Query().Get("team"),
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

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		var vars WorkflowVars

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				logger.Print("ParseForm error: " + err.Error())
				return err
			}

			assignTeam := r.FormValue("assignTeam")
			if assignTeam == "0" {
				vars.Errors = sirius.ValidationErrors{
					"selection": map[string]string{"": "Please select a team"},
				}
			}

			//this is where it picks up the new user to assign task to
			newAssigneeIdForTask, err := getAssigneeIdForTask(logger, assignTeam, r.FormValue("assignCM"))
			if err != nil {
				logger.Print("getAssigneeIdForTask error: " + err.Error())
				return err
			}

			selectedTasks := r.Form["selected-tasks"]

			// Attempt to save
			err = client.AssignTasksToCaseManager(ctx, newAssigneeIdForTask, selectedTasks)
			if err != nil {
				logger.Print("AssignTasksToCaseManager: " + err.Error())
				return err
			}

			if vars.Errors == nil {
				vars.SuccessMessage = fmt.Sprintf("%d tasks have been reassigned", len(selectedTasks))
			}
		}

		params := r.URL.Query()

		page, _ := strconv.Atoi(params.Get("page"))
		if page < 1 {
			page = 1
		}

		tasksPerPage := getTasksPerPage(params.Get("per-page"))

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

		var userSelectedAssignees []string
		if params.Has("assignee") {
			userSelectedAssignees = params["assignee"]
		}
		selectedAssignees := userSelectedAssignees
		selectedUnassigned := params.Get("unassigned")

		if selectedUnassigned == selectedTeam.Selector {
			selectedAssignees = append(selectedAssignees, strconv.Itoa(selectedTeam.Id))
			for _, t := range selectedTeam.Teams {
				selectedAssignees = append(selectedAssignees, strconv.Itoa(t.Id))
			}
		}

		var selectedTaskTypes []string
		if params.Has("task-type") {
			selectedTaskTypes = params["task-type"]
		}

		taskTypes, err := client.GetTaskTypes(ctx, selectedTaskTypes)
		if err != nil {
			logger.Print("GetTaskTypes error " + err.Error())
			return err
		}

		taskList, err := client.GetTaskList(ctx, page, tasksPerPage, selectedTeam, selectedTaskTypes, taskTypes, selectedAssignees)

		if err != nil {
			logger.Print("GetTaskList error " + err.Error())
			return err
		}
		if page > taskList.Pages.PageTotal && taskList.Pages.PageTotal > 0 {
			page = taskList.Pages.PageTotal
			taskList, err = client.GetTaskList(ctx, page, tasksPerPage, selectedTeam, selectedTaskTypes, taskTypes, selectedAssignees)

			if err != nil {
				logger.Print("GetTaskList error " + err.Error())
				return err
			}
		}

		pageDetails := client.GetPageDetails(taskList, page, tasksPerPage)

		appliedFilters := sirius.GetAppliedFilters(selectedTeam, selectedAssignees, selectedUnassigned, taskTypes)

		vars.Path = r.URL.Path
		vars.XSRFToken = ctx.XSRFToken
		vars.MyDetails = myDetails
		vars.TaskList = taskList
		vars.PageDetails = pageDetails
		vars.LoadTasks = taskTypes
		vars.TeamSelection = teamSelection
		vars.SelectedTeam = selectedTeam
		vars.SelectedAssignees = userSelectedAssignees
		vars.SelectedUnassigned = selectedUnassigned
		vars.SelectedTaskTypes = selectedTaskTypes
		vars.AppliedFilters = appliedFilters

		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
