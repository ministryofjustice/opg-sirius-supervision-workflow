package sirius

import (
	"encoding/json"
	"net/http"
)

type supervisionCaseOwnerDetail struct {
	DisplayName            string `json:"displayName"`
	SupervisionCaseOwnerId int    `json:"id"`
}

type clientDetails struct {
	CaseRecNumber        string                     `json:"caseRecNumber"`
	TaskFirstname        string                     `json:"firstname"`
	ClientId             int                        `json:"id"`
	ClientMiddlenames    string                     `json:"middlenames"`
	ClientSalutation     string                     `json:"salutation"`
	SupervisionCaseOwner supervisionCaseOwnerDetail `json:"supervisionCaseOwner"`
	TaskSurname          string                     `json:"surname"`
	ClientUId            string                     `json:"uId"`
}

type caseItemsDetails struct {
	CaseRecNumber string        `json:"caseRecNumber"`
	CaseSubtype   string        `json:"caseSubtype"`
	CaseType      string        `json:"caseType"`
	Client        clientDetails `json:"client"`
	CaseItemsId   int           `json:"id"`
	CaseItemsUId  string        `json:"uId"`
}

// type AssigneeDetails struct {
// 	DisplayName string `json:"displayName"`
// 	AssigneeId  int    `json:"id"`
// }

type ApiTask struct {
	Assignee struct {
		DisplayName string `json:"displayName"`
		Id          int    `json:"id"`
	} `json:"assignee"`
	// Assignee map[string]AssigneeDetails `json:"assignee"`
	CaseItems   []caseItemsDetails `json:"caseItems"`
	Clients     []string           `json:"clients"`
	CreatedTime string             `json:"createdTime"`
	Description string             `json:"description"`
	DueDate     string             `json:"dueDate"`
	ApiTaskId   int                `json:"id"`
	Name        string             `json:"name"`
	Persons     []string           `json:"persons"`
	RagRating   int                `json:"ragRating"`
	Status      string             `json:"status"`
	Tasktype    string             `json:"type"`
}

type TaskList struct {
	AllTaskList []ApiTask `json:"tasks"` //look into the type of this map for next time
}

// type MyTaskList struct {
// 	Name                string
// 	AssigneeDisplayName string
// 	AssigneeId          int
// }

// "tasks":[{"id":123,"type":"CWGN","status":"Not started","dueDate":"01\/10\/2019","name":"Case work - General","description":"Case w   │
// │   ork - General","ragRating":3,"assignee":{"id":65,"displayName":"case manager"},"createdTime":"25\/03\/2021 14:34:14","caseItems":[{"id":63,"uId":"7000-0000-2449","caseRecNumber":"05563462","client":{"   │
// │   id":72,"uId":"7000-0000-2357","caseRecNumber":"84029229","salutation":"Duke","firstname":"John","middlenames":"","surname":"Fearless","supervisionCaseOwner":{"id":12,"displayName":"Allocations - (Supe   │
// │   rvision)"}},"caseType":"ORDER","caseSubtype":"hw"}],"persons":[],"clients":[],"caseOwnerTask":false}]

// "tasks":[
//    {
//       "id":123,
//       "type":"CWGN",
//       "status":"Not started",
//       "dueDate":"01\/10\/2019",
//       "name":"Case work - General",
//       "description":"Case work - General",
//       "ragRating":3,
//       "assignee":{
//          "id":65,
//          "displayName":"case manager"
//       },
//       "createdTime":"25\/03\/2021 14:34:14",
//       "caseItems":[
//          {
//             "id":63,
//             "uId":"7000-0000-2449",
//             "caseRecNumber":"05563462",
//             "client":{
//                "id":72,
//                "uId":"7000-0000-2357",
//                "caseRecNumber":"84029229",
//                "salutation":"Duke",
//                "firstname":"John",
//                "middlenames":"",
//                "surname":"Fearless",
//                "supervisionCaseOwner":{
//                   "id":12,
//                   "displayName":"Allocations - (Supervision)"
//                }
//             },
//             "caseType":"ORDER",
//             "caseSubtype":"hw"
//          }
//       ],
//       "persons":[

//       ],
//       "clients":[

//       ],
//       "caseOwnerTask":false
//    }
// ]

func (c *Client) GetTaskList(ctx Context) ([]ApiTask, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/assignees/65/tasks", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// io.Copy(os.Stdout, resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v TaskList
	// var v []ApiTask
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	allTaskList := v.AllTaskList

	// allTaskList := make([]MyTaskList, len(v))

	var taskList []ApiTask

	for _, u := range allTaskList {
		task := ApiTask{
			// Assignee:    u.Assignee,
			// CaseItems:   u.CaseItems,
			Clients:     u.Clients,
			CreatedTime: u.CreatedTime,
			Description: u.Description,
			DueDate:     u.DueDate,
			ApiTaskId:   u.ApiTaskId,
			Name:        u.Name,
			Persons:     u.Persons,
			RagRating:   u.RagRating,
			Status:      u.Status,
			Tasktype:    u.Tasktype,
		}
		task.Assignee.DisplayName = u.Assignee.DisplayName
		task.Assignee.Id = u.Assignee.Id

		taskList = append(taskList, task)
	}

	// for i, t := range v {
	// 	allTaskList[i] = MyTaskList{
	// 		Name: t.Name,
	// 	}

	// if t.Assignee != nil {
	// 	allTaskList[i].AssigneeDisplayName = t.Assignee.DisplayName
	// 	allTaskList[i].AssigneeId = t.Assignee.Id
	// }
	// }

	return allTaskList, err
}
