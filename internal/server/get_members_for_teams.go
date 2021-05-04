package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type GetMembersForTeamsInformation interface {
	GetMembersForTeam(sirius.Context, int, int) (sirius.TeamSelected, error)
	SiriusUserDetails(sirius.Context) (sirius.UserDetails, error)
}

type getMembersForTeamsVars struct {
	Path         string
	XSRFToken    string
	MyDetails    sirius.UserDetails
	TeamSelected sirius.TeamSelected
	Error        string
	Errors       sirius.ValidationErrors
}

func infoForGetMembersForTeams(client WorkflowInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			fmt.Println(StatusError(http.StatusMethodNotAllowed))
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		selectedTeamToAssignTask, _ := strconv.Atoi(r.FormValue("assignTeam"))

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

		selectedTeamMembers, err := client.GetMembersForTeam(ctx, loggedInTeamId, selectedTeamToAssignTask)
		if err != nil {
			return err
		}

		vars := workflowVars{
			Path:         r.URL.Path,
			XSRFToken:    ctx.XSRFToken,
			MyDetails:    myDetails,
			TeamSelected: selectedTeamMembers,
		}

		if err != nil {
			return StatusError(http.StatusNotFound)
		}

		switch r.Method {
		case http.MethodGet:

			if err != nil {
				return err
			}

			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
