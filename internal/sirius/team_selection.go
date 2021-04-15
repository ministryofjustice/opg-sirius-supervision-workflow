package sirius

import (
	"encoding/json"
	"log"
	"net/http"
)

// {
// 	"id":21,
// 	"name":"Debt Management - (Supervision)",
// 	// "phoneNumber":"0123456789",
// 	// "displayName":"Debt Management - (Supervision)",
// 	// "deleted":false,
// 	// "email":"DebtManagement.team@opgtest.com",
// 	"members":[
// 		{"id": 79,
// 		"name":"Finance",
// 		"phoneNumber":"12345678",
// 		"displayName":"Finance User",
// 		"deleted":false,
// 		"email":"finance.user@opgtest.com"
// 		},
// 		{"id":80,
// 		"name":"Finance",
// 		"phoneNumber":"12345678",
// 		"displayName":"Finance Reporting",
// 		"deleted":false,
// 		"email":"finance.reporting@opgtest.com"
// 		}
// 		],
// 		// "children":[],
// 		// "teamType":{
// 		// 	"handle":"FINANCE",
// 		// 	"label":"Finance"
// 		// }
// }

type TeamMembers struct {
	// TeamMembersDeleted      bool   `json:"deleted"`
	// TeamMembersDisplayName  string `json:"displayName"`
	// TeamMembersEmail        string `json:"email"`
	TeamMembersId   int    `json:"id"`
	TeamMembersName string `json:"name"`
	// TeamMembersPhoneNumeber string `json:"phoneNumber"`
}

// type TeamType struct {
// 	TeamTypeDeprecated bool   `json:"deprecated"`
// 	TeamTypeHandle     string `json:"handle"`
// 	TeamTypeLabel      string `json:"label"`
// }

type TeamCollection struct {
	// Children    []string      `json:"children"`
	// Delete      bool          `json:"deleted"`
	// DisplayName string        `json:"displayName"`
	// Email       string        `json:"email"`
	// GroupName   string        `json:"groupName"`
	Id      int           `json:"id"`
	Members []TeamMembers `json:"members"`
	Name    string        `json:"name"`
	// Parent      string        `json:"parent"`
	// PhoneNumber string        `json:"phoneNumber"`
	// TeamTypeHandle TeamType      `json"teamType"`
	UserSelectedTeam int
	Kate             int
}

var selectedTeamId int

func (c *Client) GetTeamSelection(ctx Context, myDetails UserDetails, selectedTeamName int, oldTeamId int) ([]TeamCollection, error) {
	log.Println("start function oldTeamId")
	log.Println(oldTeamId)
	var v []TeamCollection

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams", nil)

	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, err
	}

	// if selectedTeamName == 0 {
	// 	selectedTeamId = myDetails.Teams[0].TeamId
	// } else {
	// 	selectedTeamId = selectedTeamName
	// }

	log.Println("team id before if statement")
	log.Println(selectedTeamId)

	if selectedTeamName == 0 && oldTeamId == 0 {
		selectedTeamId = myDetails.Teams[0].TeamId //first log on everything zero
	} else if selectedTeamName == 0 {
		selectedTeamId = oldTeamId //if submitted through page 2 take old value
	} else {
		selectedTeamId = selectedTeamName //if new value take that and add it into the struct

	}
	log.Println("team id after if statement")
	log.Println(selectedTeamId)

	for i, _ := range v {
		v[i].UserSelectedTeam = selectedTeamId
		v[i].Kate = selectedTeamId
	}

	log.Println("kate ")
	log.Println(v[0].Kate)
	// log.Println("team selection end function v")
	// log.Println(v)

	// io.Copy(os.Stdout, resp.Body)
	return v, err
}
