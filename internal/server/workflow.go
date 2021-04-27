package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type WorkflowInformation interface {
	SiriusUserDetails(sirius.Context) (sirius.UserDetails, error)
	GetTaskType(sirius.Context) (sirius.TaskTypes, error)
	GetTaskList(sirius.Context, int, int, int, int) (sirius.TaskList, sirius.TaskDetails, error)
	GetTeamSelection(sirius.Context, int, int, sirius.TeamSelected) ([]sirius.TeamCollection, error)
	GetMembersForTeam(sirius.Context, int, int) (sirius.TeamSelected, error)
	AssignTasksToCaseManager(sirius.Context, int, string) error
}

type workflowVars struct {
	Path           string
	XSRFToken      string
	MyDetails      sirius.UserDetails
	TaskList       sirius.TaskList
	TaskDetails    sirius.TaskDetails
	LoadTasks      sirius.TaskTypes
	TeamSelection  []sirius.TeamCollection
	TeamStoredData sirius.TeamStoredData
	TeamSelected   sirius.TeamSelected
	SuccessMessage string
}

type editTaskVars struct {
	Path      string
	XSRFToken string
	Success   bool
	Errors    sirius.ValidationErrors
}

func loggingInfoForWorflow(client WorkflowInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			fmt.Println(StatusError(http.StatusMethodNotAllowed))
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		search, _ := strconv.Atoi(r.FormValue("page"))
		displayTaskLimit, _ := strconv.Atoi(r.FormValue("tasksPerPage"))
		selectedTeamName, _ := strconv.Atoi(r.FormValue("change-team"))
		selectedTeamToAssignTask, _ := strconv.Atoi(r.FormValue("assignTeam"))

		myDetails, err := client.SiriusUserDetails(ctx)
		if err != nil {
			return err
		}

		loggedInTeamId := myDetails.Teams[0].TeamId

		loadTaskTypes, err := client.GetTaskType(ctx)
		if err != nil {
			return err
		}

		taskList, taskdetails, err := client.GetTaskList(ctx, search, displayTaskLimit, selectedTeamName, loggedInTeamId)
		if err != nil {
			return err
		}

		selectedTeamMembers, err := client.GetMembersForTeam(ctx, loggedInTeamId, selectedTeamToAssignTask)
		if err != nil {
			return err
		}

		teamSelection, err := client.GetTeamSelection(ctx, loggedInTeamId, selectedTeamName, selectedTeamMembers)
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
			TeamSelected:  selectedTeamMembers,
		}

		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		switch r.Method {
		case http.MethodGet:

			return tmpl.ExecuteTemplate(w, "page", vars)
		case http.MethodPost:

			checkTaskHasIdForAssigning := r.PostFormValue("assignCM")
			var newAssigneeIdForTask int

			if checkTaskHasIdForAssigning != "" {
				newAssigneeIdForTask, _ = strconv.Atoi(r.PostFormValue("assignCM"))
			} else {
				newAssigneeIdForTask = vars.TeamSelected.Id
			}

			r.ParseForm()
			taskIdArray := (r.Form["selected-tasks"])

			taskIdForUrl := ""

			for i := 0; i < len(taskIdArray); i++ {
				taskIdForUrl += taskIdArray[i]
				if i < (len(taskIdArray) - 1) {
					taskIdForUrl += "+"
				}
			}

			assignTaskVars := editTaskVars{
				Path:      r.URL.Path,
				XSRFToken: ctx.XSRFToken,
			}

			if err != nil {
				return err
			}

			// Attempt to save
			err = client.AssignTasksToCaseManager(ctx, newAssigneeIdForTask, taskIdForUrl)

			if err != nil {
				return err
			}

			if _, ok := err.(sirius.ClientError); ok {
				assignTaskVars.Errors = sirius.ValidationErrors{
					"firstname": {
						"": err.Error(),
					},
				}
				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if err != nil {
				return err
			}

			assignTaskVars.Success = true
			vars.SuccessMessage = fmt.Sprintf("%d tasks have been reassigned", len(taskIdArray))
			TaskList, _, err := client.GetTaskList(ctx, search, displayTaskLimit, selectedTeamName, loggedInTeamId)

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
