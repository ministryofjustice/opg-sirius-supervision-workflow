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
	GetTaskList(sirius.Context, int, int, int, int) (sirius.TaskList, sirius.TaskDetails, error)
	GetTeamSelection(sirius.Context, int, int, sirius.TeamSelected) ([]sirius.TeamCollection, error)
	GetMembersForTeam(sirius.Context, int, int) (sirius.TeamSelected, error)
	AssignTasksToCaseManager(sirius.Context, int, int) error
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
}

type editTaskVars struct {
	Path      string
	XSRFToken string
	// TaskId                int
	// TeamId                int
	// AssignedCaseManagerId int
	Success bool
	Errors  sirius.ValidationErrors
}

func loggingInfoForWorflow(client WorkflowInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		// add if post request bit here
		// amend following line to allow more than a get
		// if r.Method != http.MethodGet && r.Method != http.MethodPost {
		// 	return StatusError(http.StatusMethodNotAllowed)
		// }

		ctx := getContext(r)

		search, _ := strconv.Atoi(r.FormValue("page"))
		displayTaskLimit, _ := strconv.Atoi(r.FormValue("tasksPerPage"))
		selectedTeamName, _ := strconv.Atoi(r.FormValue("change-team"))
		selectedTeamToAssignTask, _ := strconv.Atoi(r.FormValue("assignTeam"))
		//selectedTask, _ := strconv.Atoi(r.FormValue("select-task"))

		myDetails, err := client.SiriusUserDetails(ctx)
		loggedInTeamId := myDetails.Teams[0].TeamId

		loadTaskTypes, err := client.GetTaskType(ctx)
		taskList, taskdetails, err := client.GetTaskList(ctx, search, displayTaskLimit, selectedTeamName, loggedInTeamId)

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

		//post methods
		//pulls id from url
		//taskId := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/"))

		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		log.Println("start of post stuff")

		switch r.Method {
		case http.MethodGet:

			return tmpl.ExecuteTemplate(w, "page", vars)
		case http.MethodPost:
			tempTeamId, err := strconv.Atoi(r.PostFormValue("assignTeam"))
			log.Println("teamId")
			log.Println(tempTeamId)
			tempAssignedCaseManagerId, err := strconv.Atoi(r.PostFormValue("assignCM"))
			log.Println("assignCM")
			log.Println(tempAssignedCaseManagerId)
			tempTaskId, err := strconv.Atoi(r.PostFormValue("selected-tasks"))
			log.Println("task id")
			log.Println(tempTaskId)

			assignTaskVars := editTaskVars{
				Path:      r.URL.Path,
				XSRFToken: ctx.XSRFToken,
				// TeamId:                tempTeamId,
				// AssignedCaseManagerId: tempAssignedCaseManagerId,
				// TaskId:                tempTaskId,
			}

			log.Println("assignTaskVars")
			log.Println(assignTaskVars)

			if err != nil {
				return err
			}

			// Attempt to save
			err = client.AssignTasksToCaseManager(ctx, tempAssignedCaseManagerId, tempTaskId)

			if err != nil {
				return err
			}
			log.Println("back after call to assign tasks")

			if _, ok := err.(sirius.ClientError); ok {
				assignTaskVars.Errors = sirius.ValidationErrors{
					"firstname": {
						"": err.Error(),
					},
				}
				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", assignTaskVars)
			}

			if err != nil {
				return err
			}

			assignTaskVars.Success = true

			return tmpl.ExecuteTemplate(w, "page", assignTaskVars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
