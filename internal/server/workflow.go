package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/logging"
	"net/http"
	"os"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"golang.org/x/sync/errgroup"
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
	TeamId         int
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

func loggingInfoForWorkflow(client WorkflowInformation, tmpl Template, defaultWorkflowTeam int) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var loggedInTeamId int
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
			logger.Print("ParseForm error: " + err.Error())
			return err
		}

		taskTypeSelected := r.Form["selected-task-type"]
		assigneeSelected := r.Form["selected-assignee"]

		data := workflowVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
		}

		//need to make two groups and make sure that the second group the first one comes back before the second group
		//think about putting the groups into methods
		//do a process tree
		//can we call the api once to get everything?
		group, groupCtx := errgroup.WithContext(ctx.Context)
		groupTwo, _ := errgroup.WithContext(ctx.Context)
		groupThree, groupCtxThree := errgroup.WithContext(ctx.Context)
		groupFour, groupCtxFour := errgroup.WithContext(ctx.Context)
		groupFive, _ := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			myDetails, err := client.GetCurrentUserDetails(ctx.With(groupCtx))
			if err != nil {
				logger.Print("GetCurrentUserDetails error " + err.Error())
				return err
			}
			data.MyDetails = myDetails
			return nil
		})

		group.Go(func() error {
			loadTaskTypes, err := client.GetTaskTypes(ctx.With(groupCtx), taskTypeSelected)
			if err != nil {
				logger.Print("GetTaskTypes error " + err.Error())
				return err
			}
			data.LoadTasks = loadTaskTypes
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		groupTwo.Go(func() error {
			loggedInTeamId = getLoggedInTeam(data.MyDetails, defaultWorkflowTeam)
			return nil
		})

		if err := groupTwo.Wait(); err != nil {
			return err
		}

		groupThree.Go(func() error {
			taskList, teamId, err := client.GetTaskList(ctx.With(groupCtxThree), search, displayTaskLimit, selectedTeamId, loggedInTeamId, taskTypeSelected, data.LoadTasks, assigneeSelected)
			if err != nil {
				logger.Print("GetTaskList error " + err.Error())
				return err
			}
			data.TeamId = teamId
			data.TaskList = taskList
			return nil
		})
		if search > data.TaskList.Pages.PageTotal && data.TaskList.Pages.PageTotal > 0 {
			groupThree.Go(func() error {
				search = data.TaskList.Pages.PageTotal
				data.TaskList, data.TeamId, err = client.GetTaskList(ctx.With(groupCtxThree), search, displayTaskLimit, selectedTeamId, loggedInTeamId, taskTypeSelected, data.LoadTasks, assigneeSelected)
				if err != nil {
					logger.Print("GetTaskList error " + err.Error())
					return err
				}
				return nil
			})
		}

		if err := groupThree.Wait(); err != nil {
			return err
		}

		groupFour.Go(func() error {
			pageDetails := client.GetPageDetails(data.TaskList, search, displayTaskLimit)
			data.PageDetails = pageDetails
			return nil
		})

		groupFour.Go(func() error {
			teamSelection, err := client.GetTeamsForSelection(ctx.With(groupCtxFour), data.TeamId, assigneeSelected)
			if err != nil {
				logger.Print("GetTeamsForSelection error " + err.Error())
				return err
			}
			data.TeamSelection = teamSelection
			return nil
		})

		groupFour.Go(func() error {
			assigneesForFilter, err := client.GetAssigneesForFilter(ctx.With(groupCtxFour), data.TeamId, assigneeSelected)
			if err != nil {
				logger.Print("GetAssigneesForFilter error " + err.Error())
				return err
			}

			data.Assignees = assigneesForFilter
			return nil
		})

		if err := groupFour.Wait(); err != nil {
			return err
		}

		groupFive.Go(func() error {
			appliedFilters := client.GetAppliedFilters(data.TeamId, data.LoadTasks, data.TeamSelection, data.Assignees)
			data.AppliedFilters = appliedFilters
			return nil
		})

		if err := groupFive.Wait(); err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", data)
		case http.MethodPost:
			var newAssigneeIdForTask int
			selectedTeamToAssignTaskString := r.FormValue("assignTeam")
			if selectedTeamToAssignTaskString == "0" {
				data.Errors = sirius.ValidationErrors{
					"selection": map[string]string{"": "Please select a team"},
				}

				return tmpl.ExecuteTemplate(w, "page", data)
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

			taskIdArray := (r.Form["selected-tasks"])
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

			if data.Errors == nil {
				data.SuccessMessage = fmt.Sprintf("%d tasks have been reassigned", len(taskIdArray))
			}

			data.TaskList, _, err = client.GetTaskList(ctx, search, displayTaskLimit, selectedTeamId, loggedInTeamId, taskTypeSelected, data.LoadTasks, assigneeSelected)
			if err != nil {
				logger.Print("data.TaskList error: " + err.Error())
				return err
			}

			data.PageDetails = client.GetPageDetails(data.TaskList, search, displayTaskLimit)

			return tmpl.ExecuteTemplate(w, "page", data)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
