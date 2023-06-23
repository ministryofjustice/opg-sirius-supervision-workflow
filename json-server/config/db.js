module.exports = function () {
  let repeatableTasks = [
    {
      "assignee": {
        "displayName": "Allocations User3",
        "id": 86
      },
      "caseItems": [
        {
          "client": {
            "caseRecNumber": "51132787",
            "firstname": "Lizzo",
            "id": 3333,
            "supervisionCaseOwner": {
              "displayName": "Allocations - (Supervision)"
            },
            "surname": "Surname"
          }
        }
      ],
      "dueDate": "01/02/2021",
      "name": "Case work - Complaint review"
    },
    {
      "assignee": {
        "displayName": "Allocations User3",
        "id": 69
      },
      "caseItems": [
        {
          "client": {
            "caseRecNumber": "48806317",
            "firstname": "Jimi",
            "id": 2564,
            "supervisionCaseOwner": {
              "displayName": "Allocations - (Supervision)"
            },
            "surname": "Hendrix"
          }
        }
      ],
      "dueDate": "28/02/2018",
      "name": "Order - Allocate to team"
    },
    {
      "assignee": {
        "displayName": "Margaret Bavaria-Straubing",
        "id": 99
      },
      "caseItems": [
        {
          "client": {
            "caseRecNumber": "83948321",
            "firstname": "Elton",
            "id": 1458,
            "supervisionCaseOwner": {
              "displayName": "Allocations - (Supervision)"
            },
            "surname": "John"
          }
        }
      ],
      "dueDate": "01/12/2018",
      "name": "Case work - Call back request"
    },
    {
      "assignee": {
        "displayName": "Allocations User13",
        "id": 88
      },
      "caseItems": [
        {
          "client": {
            "caseRecNumber": "58748192",
            "firstname": "Elvis",
            "id": 1237,
            "supervisionCaseOwner": {
              "displayName": "Allocations - (Supervision)"
            },
            "surname": "Presley"
          }
        }
      ],
      "dueDate": "01/02/2021",
      "name": "Case work - Complaint review"
    },
    {
      "assignee": {
        "displayName": "Allocations User10",
        "id": 12
      },
      "caseItems": [
        {
          "client": {
            "caseRecNumber": "33731087",
            "firstname": "Sean",
            "id": 4561,
            "supervisionCaseOwner": {
              "displayName": "Allocations - (Supervision)"
            },
            "surname": "Bean"
          }
        }
      ],
      "dueDate": "01/02/2021",
      "name": "Case work - Safeguarding"
    },
    {
      "assignee": {
        "displayName": "Allocations - (Supervision)",
        "id": 58
      },
      "caseItems": [
        {
          "client": {
            "caseRecNumber": "91221815",
            "firstname": "David",
            "id": 4561,
            "supervisionCaseOwner": {
              "displayName": "Allocations - (Supervision)"
            },
            "surname": "Attenborough"
          }
        }
      ],
      "dueDate": "01/02/2021",
      "name": "Correspondence - Review failed draft"
    },
    {
      "assignee": {
        "displayName": "Allocations - (Supervision)",
        "id": 94
      },
      "caseItems": [
        {
          "client": {
            "caseRecNumber": "52556905",
            "firstname": "Alexander",
            "id": 7894,
            "supervisionCaseOwner": {
              "displayName": "Allocations - (Supervision)"
            },
            "surname": "Bell"
          }
        }
      ],
      "dueDate": "01/02/2021",
      "name": "Finance - Sop processing"
    },
    {
      "assignee": {
        "displayName": "Allocations User3",
        "id": 77
      },
      "caseItems": [
        {
          "client": {
            "caseRecNumber": "94277087",
            "firstname": "Henry",
            "id": 3691,
            "supervisionCaseOwner": {
              "displayName": "Allocations - (Supervision)"
            },
            "surname": "Cavill"
          }
        }
      ],
      "dueDate": "01/02/2021",
      "name": "Report - Full staff review"
    },
    {
      "assignee": {
        "displayName": "Allocations User1",
        "id": 46
      },
      "caseItems": [
        {
          "client": {
            "caseRecNumber": "93766765",
            "firstname": "Henry",
            "id": 2583,
            "supervisionCaseOwner": {
              "displayName": "Allocations - (Supervision)"
            },
            "surname": "Ford"
          }
        }
      ],
      "dueDate": "01/02/2021",
      "name": "Visit - Review report"
    },
    {
      "assignee": {
        "displayName": "Allocations - (Supervision)",
        "id": 64
      },
      "caseItems": [
        {
          "client": {
            "caseRecNumber": "62168380",
            "firstname": "Matt",
            "id": 1472,
            "supervisionCaseOwner": {
              "displayName": "Allocations - (Supervision)"
            },
            "surname": "Smith"
          }
        }
      ],
      "dueDate": "01/02/2021",
      "name": "Order - Review the order"
    }
  ]

  // Generate a number of tasks by repeating the repeatableTasks
  let generateTasks = function (number) {
    let tasks = []
    for (let i = 1; i <= number; i++) {
      let mod = i % repeatableTasks.length
      let taskNum = mod === 0 ? repeatableTasks.length : mod
      let task = structuredClone(repeatableTasks[taskNum - 1])
      task["id"] = i;
      tasks.push(task)
    }
    return tasks
  }

  return {
    "successes": [
      {
        "id": 999,
        "type":"ORAL",
        "status":"Not started",
        "dueDate":"25\/05\/2023",
        "name":"",
        "description":"A client has been created",
        "ragRating":2,
        "assignee":{
          "id":25,
          "teams":[],
          "displayName":"Lay Team 1 - (Supervision)"
        },
        "createdTime":"10\/05\/2023 09:27:05",
        "caseItems":[],
        "persons":[
          {
            "id":42,
            "uId":"7000-0000-0856",
            "caseRecNumber":"71115167",
            "salutation":"Mr",
            "firstname":"Luke",
            "middlenames":"",
            "surname":"Crete",
            "supervisionCaseOwner":{
              "id":25,
              "teams":[],
              "displayName":"Lay Team 1 - (Supervision)"}
          }
        ],
        "clients":[
          {
            "id":42,
            "uId":"7000-0000-0856",
            "caseRecNumber":"71115167",
            "salutation":"Mr",
            "firstname":"Luke",
            "middlenames":"",
            "surname":"Crete",
            "supervisionCaseOwner":
                {
                  "id":25,
                  "teams":[],
                  "displayName":"Lay Team 1 - (Supervision)"
                }
          }
        ],
        "caseOwnerTask":false,
        "isPriority":true
      }
    ],
    "assign-tasks-to-casemanager": {},
    "teams": [
      {
        "id": 13,
        "name": "Allocations - (Supervision)",
        "phoneNumber": "0123456789",
        "displayName": "Allocations - (Supervision)",
        "deleted": false,
        "email": "allocations.team@opgtest.com",
        "members": [
          {
            "id": 72,
            "name": "Allocations",
            "phoneNumber": "12345678",
            "displayName": "Allocations User1",
            "deleted": false,
            "email": "allocations@opgtest.com"
          }
        ],
        "teamType": {
          "handle": "ALLOCATIONS",
          "label": "Allocations",
          "deprecated": null
        }
      },
      {
        "id": 32,
        "name": "Investigations - (Supervision)",
        "phoneNumber": "0123456789",
        "displayName": "Investigations - (Supervision)",
        "deleted": false,
        "email": "Investigations.team@opgtest.com",
        "members": [
          {
            "id": 101,
            "name": "Investigations",
            "phoneNumber": "12345678",
            "displayName": "Investigations User1",
            "deleted": false,
            "email": "investigations@opgtest.com"
          }
        ],
        "teamType": {
          "handle": "INVESTIGATIONS",
          "label": "Investigations",
          "deprecated": null
        }
      },
      {
        "id": 21,
        "name": "Lay Team 1 - (Supervision)",
        "phoneNumber": "0123456789",
        "displayName": "Lay Team 1 - (Supervision)",
        "deleted": false,
        "email": "Allocations.team@opgtest.com",
        "members": [
          {
            "id": 766,
            "name": "LayTeam1 User1",
            "phoneNumber": "12345678",
            "displayName": "LayTeam1 User1",
            "deleted": false,
            "email": "lay1-user1@opgtest.com"
          },
        ],
        "teamType": {
          "handle": "LAY",
          "label": "Lay",
          "deprecated": null
        }
      },
      {
        "id": 22,
        "name": "Lay Team 2 - (Supervision)",
        "phoneNumber": "0123456789",
        "displayName": "Lay Team 2 - (Supervision)",
        "deleted": false,
        "email": "LayTeam2.team@opgtest.com",
        "members": [
          {
            "id": 93,
            "name": "LayTeam2",
            "phoneNumber": "12345678",
            "displayName": "LayTeam2 User1",
            "deleted": false,
            "email": "lay2@opgtest.com"
          }
        ],
        "teamType": {
          "handle": "LAY",
          "label": "Lay",
          "deprecated": null
        }
      },
      {
        "id": 24,
        "name": "PA Team 1 - (Supervision)",
        "phoneNumber": "0123456789",
        "displayName": "PA Team 1 - (Supervision)",
        "deleted": false,
        "email": "PATeam1.team@opgtest.com",
        "members": [
          {
            "id": 94,
            "name": "PATeam1",
            "phoneNumber": "12345678",
            "displayName": "PATeam1 User1",
            "deleted": false,
            "email": "pa1@opgtest.com"
          }
        ],
        "teamType": {
          "handle": "PA",
          "label": "PA",
          "deprecated": null
        }
      },
      {
        "id": 27,
        "name": "Pro Team 1 - (Supervision)",
        "phoneNumber": "0123456789",
        "displayName": "Pro Team 1 - (Supervision)",
        "deleted": false,
        "email": "ProTeam1.team@opgtest.com",
        "members": [
          {
            "id": 96,
            "name": "PROTeam1",
            "phoneNumber": "12345678",
            "displayName": "PROTeam1 User1",
            "deleted": false,
            "email": "pro1@opgtest.com"
          }
        ],
        "teamType": {
          "handle": "PRO",
          "label": "Pro",
          "deprecated": null
        }
      }
    ],
    "tasktypes-supervision": {
      "task_types": {
        "CWGN": {
          "category": "supervision",
          "complete": "Casework - General",
          "handle": "CWGN",
          "incomplete": "Casework - General",
          "user": true,
          "ecmTask": true
        },
        "ORAL": {
          "category": "supervision",
          "complete": "Order - Allocate to team",
          "handle": "ORAL",
          "incomplete": "Order - Allocate to team",
          "user": true,
          "ecmTask": false
        }
      }
    },
    "users-current": {
      "deleted": false,
      "displayName": "case manager",
      "email": "case.manager@opgtest.com",
      "firstname": "case",
      "id": 65,
      "locked": false,
      "name": "case",
      "phoneNumber": "12345678",
      "roles": [
        "Case Manager",
        "Manager"
      ],
      "surname": "manager",
      "suspended": false,
      "teams": [
        {
          "displayName": "Lay Team 1 - (Supervision)",
          "id": 13
        }
      ]
    },
    "tasks": [
      {
        "id": 13,
        "pages": {
          "current": 1,
          "total": 1
        },
        "tasks": generateTasks(10),
        "total": 10
      },
      {
        "id": 21,
        "pages": {
          "current": 1,
          "total": 1
        },
        "tasks": [],
        "total": 0
      },
      {
        "id": 22,
        "pages": {
          "current": 1,
          "total": 1
        },
        "tasks": [
          {
            "assignee": {
              "displayName": "LayTeam2 User4",
              "id": 86
            },
            "caseItems": [
              {
                "client": {
                  "caseRecNumber": "51132787",
                  "firstname": "Beth",
                  "id": 6354,
                  "supervisionCaseOwner": {
                    "displayName": "Lay Team 2 - (Supervision)"
                  },
                  "surname": "Lay"
                }
              }
            ],
            "dueDate": "11/11/2011",
            "name": "Case work - Complaint review"
          }
        ],
        "total": 1
      },
      {
        "id": 24,
        "pages": {
          "current": 1,
          "total": 1
        },
        "tasks": [
          {
            "assignee": {
              "displayName": "PATeam1 User1",
              "id": 86
            },
            "caseItems": [
              {
                "client": {
                  "caseRecNumber": "51132787",
                  "firstname": "Beth",
                  "id": 6354,
                  "supervisionCaseOwner": {
                    "displayName": "Lay Team 2 - (Supervision)"
                  },
                  "feePayer": {
                    "id": 12,
                    "displayName": "Mr Fee-paying Deputy",
                    "deputyType": {
                      "handle": "PA",
                      "label": "PA"
                    }
                  },
                  "surname": "Lay"
                }
              }
            ],
            "dueDate": "11/11/2011",
            "name": "Case work - Complaint review"
          }
        ],
        "total": 1
      }
    ],
    "tasks-view-25": [
      {
        "id": 1,
        "page": 1,
        "pages": {
          "current": 1,
          "total": 5
        },
        "tasks": generateTasks(25),
        "total": 101
      },
      {
        "id": 2,
        "page": 2,
        "pages": {
          "current": 2,
          "total": 5
        },
        "tasks": generateTasks(25),
        "total": 101
      },
      {
        "id": 3,
        "page": 3,
        "pages": {
          "current": 3,
          "total": 5
        },
        "tasks": generateTasks(25),
        "total": 101
      },
      {
        "id": 4,
        "page": 4,
        "pages": {
          "current": 4,
          "total": 5
        },
        "tasks": generateTasks(25),
        "total": 101
      },
      {
        "id": 5,
        "page": 5,
        "pages": {
          "current": 5,
          "total": 5
        },
        "tasks": generateTasks(1),
        "total": 101
      }
    ],
    "tasks-view-50": [
      {
        "id": 1,
        "page": 1,
        "pages": {
          "current": 1,
          "total": 2
        },
        "tasks": generateTasks(50),
        "total": 51
      },
      {
        "id": 2,
        "page": 2,
        "pages": {
          "current": 2,
          "total": 2
        },
        "tasks": generateTasks(1),
        "total": 51
      }
    ],
    "tasks-view-100": [
      {
        "id": 1,
        "page": 1,
        "pages": {
          "current": 1,
          "total": 1
        },
        "tasks": generateTasks(10),
        "total": 10
      }
    ],
    "caseload":  {
        "id": 44,
        "pages": {
          "current": 1,
          "total": 1
        },
        "total": 1,
        "clients": [
          {
            "id": 63,
            "caseRecNumber": "42687883",
            "firstname": "Ro",
            "surname": "Bot",
            "supervisionCaseOwner":{
              "id":25,
              "teams":[],
              "displayName":"Lay Team 1 - (Supervision)"
            },
            "cases": [
              {
                "id": 92,
                "caseRecNumber": "33594483",
                "latestAnnualReport": {
                  "dueDate": "21/12/2023"
                },
                "orderStatus": {
                  "handle": "CLOSED",
                  "label": "Closed",
                  "deprecated": false
                },
              }
            ],
            "supervisionLevel": "Minimal",
          },
        ]
    },
  }
}
