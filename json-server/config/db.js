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
      "type":"CWGN",
      "name": "Case work - General",
      "isPriority": true
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
      "type":"ORAL",
      "name": "Order - Allocate to team"
    },
    {
      "assignee": {
        "displayName": "Allocations - (Supervision)",
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
      "type":"CCBR",
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
      "type":"CCR",
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
      "type":"CS",
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
      "type":"CDFC",
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
      "type":"FSPR",
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
      "type":"RFSR",
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
      "type":"RVRR",
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
      "type":"ORRE",
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
          },
          {
              "id": 73,
              "name": "Allocations",
              "phoneNumber": "12345678",
              "displayName": "Allocations User2",
              "deleted": false,
              "email": "allocations2@opgtest.com"
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
          },
          {
              "id":97,
              "name": "PROTeam1",
              "phoneNumber": "12345678",
              "displayName": "PROTeam1 User2",
              "deleted": false,
              "email": "pro2@opgtest.com"
          }
        ],
        "teamType": {
          "handle": "PRO",
          "label": "Pro",
          "deprecated": null
        }
      },
      {
        "id": 28,
        "name": "Lay Team - New Deputy Orders",
        "displayName": "Lay Team - New Deputy Orders",
        "deleted": false,
        "teamType": {
          "handle": "LAY",
          "label": "Lay",
          "deprecated": null
        }
      },
      {
        "id": 29,
        "name": "Health and Welfare Team",
        "displayName": "Health and Welfare Team",
        "deleted": false,
        "teamType": {
          "handle": "HW",
          "label": "HW",
          "deprecated": null
        }
      },
      {
        "id": 0,
        "name": "Lay Deputy Team",
        "displayName": "Lay Deputy Team",
        "members": [],
      },
      {
        "id": 40,
        "name": "Supervision Closed Cases",
        "phoneNumber": "0123456789",
        "displayName": "Supervision Closed Cases",
        "deleted": false,
        "email": "ClosedCases.team@opgtest.com",
        "members": [],
        "teamType": {
          "handle": "LAY",
          "label": "Lay",
          "deprecated": null
        }
      },
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
        "CCR": {
          "category": "supervision",
          "complete": "Casework - Complaint reviewed",
          "handle": "CCR",
          "incomplete": "Casework - Complaint review",
          "user": true,
          "ecmTask": true
        },
        "CCBR": {
          "category": "supervision",
          "complete": "Casework - Call back request completed",
          "handle": "CCBR",
          "incomplete": "Casework - Call back request",
          "user": true,
          "ecmTask": true
        },
        "CS": {
          "category": "supervision",
          "complete": "Casework - Safeguarding completed",
          "handle": "CS",
          "incomplete": "Casework - Safeguarding",
          "user": true,
          "ecmTask": true
        },
        "CDFC": {
          "category": "supervision",
          "complete": "Correspondence - Review failed draft",
          "handle": "CDFC",
          "incomplete": "Correspondence - Reviewed draft failure",
          "user": true,
          "ecmTask": true
        },
        "FSPR": {
          "category": "supervision",
          "complete": "Finance - Sop processed",
          "handle": "FSPR",
          "incomplete": "Finance - Sop processing",
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
        },
        "ORRE": {
          "category": "supervision",
          "complete": "Order - Reviewed the order",
          "handle": "ORRE",
          "incomplete": "Order - Review the order",
          "user": true,
          "ecmTask": false
        },
       "RFSR": {
          "category": "supervision",
          "complete": "Report - Full staff reviewed",
          "handle": "RFSR",
          "incomplete": "Report - Full staff review",
          "user": true,
          "ecmTask": true
        },
       "RVRR": {
          "category": "supervision",
          "complete": "Visit - Reviewed report",
          "handle": "RVRR",
          "incomplete": "Visit - Review report",
          "user": true,
          "ecmTask": true
        }
      }
    },
    "tasktypes-deputy": {
      "task_types": {
        "PFU": {
          "category": "deputy",
          "complete": "PDR follow up",
          "handle": "PFU",
          "incomplete": "PDR follow up",
          "user": true,
          "proDeputyTask": true
        },
        "QCUC": {
          "category": "deputy",
          "complete": "Quarterly catch up call",
          "handle": "QCUC",
          "incomplete": "Quarterly catch up call",
          "user": true,
          "paDeputyTask": true
        },
        "AVR": {
          "category": "deputy",
          "complete": "Assurance visit review",
          "handle": "AVR",
          "incomplete": "Assurance visit review",
          "user": true,
          "paDeputyTask": true
        },
        "AVC": {
          "category": "deputy",
          "complete": "Assurance visit commission",
          "handle": "AVC",
          "incomplete": "Assurance visit commission",
          "user": true,
          "paDeputyTask": true
        },
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
        "Manager",
        "Finance Reporting"
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
        "page": 1,
        "metadata":{
            "taskTypeCount":[
                {"type":"CWGN","count":2},
                {"type":"CCR","count":2},
                {"type":"CCBR","count":2},
                {"type":"CS","count":2},
                {"type":"CDFC","count":2},
                {"type":"FSPR","count":2},
                {"type":"ORAL","count":2},
                {"type":"ORRE","count":2},
                {"type":"RFSR","count":2},
                {"type":"RVRR","count":2}
            ],
            "assigneeTaskCount":[
                {"assignee":72,"count":12},
                {"assignee":73,"count":8}
            ]
        },
        "pages": {
          "current": 1,
          "total": 1
        },
        "tasks": generateTasks(20),
        "total": 20
      },
      {
        "id": 21,
        "pages": {
          "current": 1,
          "total": 1
        },
        "metadata":{
            "taskTypeCount":[],
            "assigneeTaskCount":[]
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
        "metadata":{
            "taskTypeCount":[
                {"type":"CCR","count":1}
            ],
            "assigneeTaskCount":[
                {"assignee":86,"count":1}
            ]
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
            "type": "CCR",
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
        "metadata":{
            "taskTypeCount":[
                {"type":"CCR","count":1}
            ],
            "assigneeTaskCount":[
                {"assignee":86,"count":1}
            ]
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
            "type": "CCR",
            "name": "Case work - Complaint review"
          }
        ],
        "total": 1
      },
      {
        "id": 27,
        "pages": {
          "current": 1,
          "total": 1
        },
        "metadata":{
            "taskTypeCount":[],
            "assigneeTaskCount":[]
        },
        "tasks": [],
        "total": 0
      },
      {
        "id": 29,
        "pages": {
          "current": 1,
          "total": 1
        },
        "metadata":{
            "taskTypeCount":[],
            "assigneeTaskCount":[]
        },
        "tasks": [],
        "total": 0
      }
    ],
    "deputy-tasks": [
      {
        "id": 24,
        "pages": {
          "current": 1,
          "total": 1
        },
        "metadata":{
            "taskTypeCount":[],
            "assigneeTaskCount":[]
        },
        "tasks": [],
        "total": 0
      },
      {
        "id": 27,
        "pages": {
          "current": 1,
          "total": 1
        },
        "metadata":{
            "taskTypeCount":[
                {"type":"PFU","count":1},
                {"type":"AVC","count":1},
                {"type":"AVR","count":1},
                {"type":"QCUC","count":1}
            ],
            "assigneeTaskCount":[
                {"assignee":96,"count":3},
                {"assignee":27,"count":1}
            ]
        },
        "tasks": [
          {
            "id": 1,
            "assignee": {
              "displayName": "PROTeam1 User1",
              "id": 96
            },
            "deputies": [
              {
                "id": 13,
                "deputyNumber": 123456,
                "displayName": "Mr Fee-paying Deputy",
                "deputyType": {
                  "handle": "PRO",
                  "label": "Professional"
                },
                "deputyAddress": {
                  "town": "Derby"
                },
                "firm": {
                  "id": 78745361,
                  "firmName": "Krusty Krabs",
                  "firmNumber": 789456123
                },
              }
            ],
            "dueDate": "11/02/2021",
            "name": "PDR follow up",
            "isPriority": true
          },
          {
            "id": 2,
            "assignee": {
              "displayName": "PROTeam1 User1",
              "id": 96
            },
            "deputies": [
              {
                "id": 14,
                "deputyNumber": 654321,
                "displayName": "Mr PRO Deputy",
                "deputyType": {
                  "handle": "PRO",
                  "label": "Professional"
                },
                "deputyAddress": {
                  "town": "Nottingham"
                }
              }
            ],
            "dueDate": "01/02/2021",
            "name": "Assurance visit report due"
          },
         {
              "id": 3,
              "assignee": {
                "displayName": "PROTeam1 User1",
                "id": 96
              },
              "deputies": [
                {
                  "id": 14,
                  "deputyNumber": 654321,
                  "displayName": "Mr PRO Deputy",
                  "deputyType": {
                    "handle": "PRO",
                    "label": "Professional"
                  },
                  "deputyAddress": {
                    "town": "Nottingham"
                  }
                }
              ],
              "dueDate": "03/02/2021",
              "name": "Assurance visit review"
          },
          {
              "id": 4,
              "assignee":{"id":27,"teams":[],"displayName":"Pro Team 1 - (Supervision)"},
              "deputies": [
                {
                  "id": 14,
                  "deputyNumber": 654321,
                  "displayName": "Mr PRO Deputy",
                  "deputyType": {
                    "handle": "PRO",
                    "label": "Professional"
                  },
                  "deputyAddress": {
                    "town": "Nottingham"
                  }
                }
              ],
              "dueDate": "01/04/2021",
              "name": "Quarterly catch up call"
            }
        ],
        "total": 2
      },
    ],
    "deputies": [
      {
        "id": 27,
        "pages": {
          "current": 1,
          "total": 1
        },
        "metadata":{
            "ecmCount": [
                {"assignee" : 96, "count": 1},
                {"assignee" : 97, "count": 1}
            ],
        },
        "persons": [
          {
            "id": 13,
            "deputyNumber": 123456,
            "displayName": "Mr Fee-paying Deputy",
            "deputyType": {
              "handle": "PRO",
              "label": "Professional"
            },
            "deputyAddress": {
              "town": "Derby"
            },
            "executiveCaseManager": {
              "displayName": "PROTeam1 User1",
              "id": 96
            },
            "mostRecentlyCompletedAssurance": {
              "reportReviewDate" : "2023-05-26T00:00:00+00:00",
              "reportMarkedAs": {
                "handle": "GREEN",
                "label": "Green"
              },
              "assuranceType": {
                "handle": "VISIT",
                "label": "Visit"
              }
            },
            "activeClientCount": 100,
            "activeNonCompliantClientCount": 10,
            "deputyImportantInformation": {
              "panelDeputy": true
            },
            "firm": {
                "id": 78745361,
                "firmName": "Krusty Krabs",
                "firmNumber": 789456123
            }
          },
          {
              "id": 15,
              "deputyNumber": 5535353,
              "displayName": "Mr Horus Lupercal",
              "deputyType": {
                "handle": "PRO",
                "label": "Professional"
              },
              "deputyAddress": {
                "town": "Istavan"
              },
              "executiveCaseManager": {
                "displayName": "PROTeam1 User2",
                "id": 97
              },
              "mostRecentlyCompletedAssurance": {
                "reportReviewDate" : "2023-09-12T00:00:00+00:00",
                "reportMarkedAs": {
                  "handle": "RED",
                  "label": "Red"
                },
                "assuranceType": {
                  "handle": "ASSURANCE",
                  "label": "Assurance"
                }
              },
              "activeClientCount": 77,
              "activeNonCompliantClientCount": 4,
              "deputyImportantInformation": {},
              "firm": {
                  "id": 444455555,
                  "firmName": "Chaos Inc",
                  "firmNumber": 123456
              }
          },
        ],
        "total": 2
      },
      {
          "id": 24,
          "pages": {
            "current": 1,
            "total": 1
          },
          "metadata":{
              "ecmCount": [
                  {"assignee" : 27, "count": 1},
              ],
          },
          "persons": [
            {
               "id": 14,
              "deputyNumber": 654321,
              "displayName": "Mr PA Deputy",
              "deputyType": {
                "handle": "PA",
                "label": "Public authority"
              },
              "deputyAddress": {
                "town": "Derby"
              },
              "executiveCaseManager": {
                "id": 27,
                "displayName": "PA Team 1 - (Supervision)",
              },
              "mostRecentlyCompletedAssurance": {},
              "activeClientCount": 81,
              "activeNonCompliantClientCount": 34
            },
          ],
          "total": 1
      },
    ],
    "tasks-view-25": [
      {
        "id": 1,
        "page": 1,
        "pages": {
          "current": 1,
          "total": 5
        },
        "metadata":{
            "taskTypeCount":[
                {"type":"CWGN","count":3},
                {"type":"CCR","count":3},
                {"type":"CCBR","count":3},
                {"type":"CS","count":3},
                {"type":"CDFC","count":3},
                {"type":"FSPR","count":2},
                {"type":"ORAL","count":2},
                {"type":"ORRE","count":2},
                {"type":"RFSR","count":2},
                {"type":"RVRR","count":2}
            ],
            "assigneeTaskCount":[
              {"assignee":72,"count":12},
              {"assignee":73,"count":8}
            ]
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
        "metadata":{
            "taskTypeCount":[
                {"type":"CWGN","count":3},
                {"type":"CCR","count":3},
                {"type":"CCBR","count":3},
                {"type":"CS","count":3},
                {"type":"CDFC","count":3},
                {"type":"FSPR","count":2},
                {"type":"ORAL","count":2},
                {"type":"ORRE","count":2},
                {"type":"RFSR","count":2},
                {"type":"RVRR","count":2}
            ],
             "assigneeTaskCount":[
                   {"assignee":72,"count":12},
                   {"assignee":73,"count":8}
            ]
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
        "metadata":{
            "taskTypeCount":[
                {"type":"CWGN","count":3},
                {"type":"CCR","count":3},
                {"type":"CCBR","count":3},
                {"type":"CS","count":3},
                {"type":"CDFC","count":3},
                {"type":"FSPR","count":2},
                {"type":"ORAL","count":2},
                {"type":"ORRE","count":2},
                {"type":"RFSR","count":2},
                {"type":"RVRR","count":2}
            ],
            "assigneeTaskCount":[
              {"assignee":72,"count":12},
              {"assignee":73,"count":8}
            ]
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
        "metadata":{
            "taskTypeCount":[
                {"type":"CWGN","count":3},
                {"type":"CCR","count":3},
                {"type":"CCBR","count":3},
                {"type":"CS","count":3},
                {"type":"CDFC","count":3},
                {"type":"FSPR","count":2},
                {"type":"ORAL","count":2},
                {"type":"ORRE","count":2},
                {"type":"RFSR","count":2},
                {"type":"RVRR","count":2}
            ],
            "assigneeTaskCount":[
              {"assignee":72,"count":12},
              {"assignee":73,"count":8}
            ]
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
        "metadata":{
            "taskTypeCount":[
                {"type":"CWGN","count":3},
                {"type":"CCR","count":3},
                {"type":"CCBR","count":3},
                {"type":"CS","count":3},
                {"type":"CDFC","count":3},
                {"type":"FSPR","count":2},
                {"type":"ORAL","count":2},
                {"type":"ORRE","count":2},
                {"type":"RFSR","count":2},
                {"type":"RVRR","count":2}
            ],
            "assigneeTaskCount":[
              {"assignee":72,"count":12},
              {"assignee":73,"count":8}
            ]
        },
        "tasks": generateTasks(25),
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
        "metadata":{
            "taskTypeCount":[
                {"type":"CWGN","count":5},
                {"type":"CCR","count":5},
                {"type":"CCBR","count":5},
                {"type":"CS","count":5},
                {"type":"CDFC","count":5},
                {"type":"FSPR","count":5},
                {"type":"ORAL","count":5},
                {"type":"ORRE","count":5},
                {"type":"RFSR","count":5},
                {"type":"RVRR","count":5}
            ],
            "assigneeTaskCount":[
              {"assignee":72,"count":12},
              {"assignee":73,"count":8}
            ]
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
        "metadata":{
            "taskTypeCount":[
                {"type":"CWGN","count":6},
                {"type":"CCR","count":5},
                {"type":"CCBR","count":5},
                {"type":"CS","count":5},
                {"type":"CDFC","count":5},
                {"type":"FSPR","count":5},
                {"type":"ORAL","count":5},
                {"type":"ORRE","count":5},
                {"type":"RFSR","count":5},
                {"type":"RVRR","count":5}
            ],
            "assigneeTaskCount":[
               {"assignee":72,"count":12},
               {"assignee":73,"count":8}
            ]
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
        "metadata":{
            "taskTypeCount":[
                {"type":"CWGN","count":1},
                {"type":"CCR","count":1},
                {"type":"CCBR","count":1},
                {"type":"CS","count":1},
                {"type":"CDFC","count":1},
                {"type":"FSPR","count":1},
                {"type":"ORAL","count":1},
                {"type":"ORRE","count":1},
                {"type":"RFSR","count":1},
                {"type":"RVRR","count":1}
            ],
            "assigneeTaskCount":[
               {"assignee":72,"count":12},
               {"assignee":73,"count":8}
            ]
        },
        "tasks": generateTasks(10),
        "total": 10
      }
    ],
    "caseload":  [
      {
        "id": 21,
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
            "supervisionCaseOwner": {
              "id": 21,
              "teams": [],
              "displayName": "Lay Team 1 - (Supervision)"
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
            "supervisionLevel": {
              "handle": "MINIMAL",
              "label": "Minimal"
            }
          },
        ]
      },
      {
        "id": 28,
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
            "supervisionCaseOwner": {
              "id": 28,
              "teams": [],
              "displayName": "Lay Team - New Deputy Orders"
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
                "orderDate": "01/01/2020",
                "madeActiveDate": "02/01/2020",
                "howDeputyAppointed": {
                  "handle": "SOLE",
                  "label": "Sole",
                  "deprecated": false
                },
                "introductoryTargetDate": "21/02/2020"
              }
            ],
            "supervisionLevel": {
              "handle": "MINIMAL",
              "label": "Minimal"
            }
          },
        ]
      },
      {
        "id": 29,
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
            "supervisionCaseOwner": {
              "id": 29,
              "teams": [],
              "displayName": "Health & Welfare Team"
            },
            "activeCaseType": {
              "handle": "HYBRID",
              "label": "Hybrid",
            },
            "cases": [
              {
                "id": 92,
                "caseRecNumber": "33594483",
                "latestAnnualReport": {
                  "dueDate": "21/12/2023"
                },
                "caseSubtype": "hw",
                "orderStatus": {
                  "handle": "ACTIVE",
                  "label": "Active",
                  "deprecated": false
                },
                "orderDate": "01/01/2020",
                "madeActiveDate": "02/01/2020",
                "howDeputyAppointed": {
                  "handle": "SOLE",
                  "label": "Sole",
                  "deprecated": false
                },
                "introductoryTargetDate": "21/02/2020"
              }
            ],
            "supervisionLevel": {
              "handle": "MINIMAL",
              "label": "Minimal"
            },
            "DeputyTypes":
            [
              {
                "handle": "PA",
                "label": "PA"
              },
            ],
          },
        ]
      },
      {
        "id": 40,
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
            "supervisionCaseOwner": {
              "id": 40,
              "teams": [],
              "displayName": "Supervision Closed Cases"
            },
            "cases": [
              {
                "id": 92,
                "caseRecNumber": "42687883",
                "orderStatus": {
                  "handle": "CLOSED",
                  "label": "Closed",
                  "deprecated": false
                },
              },
              {
                "id": 93,
                "caseRecNumber": "42687883",
                "orderStatus": {
                  "handle": "CLOSED",
                  "label": "Closed",
                  "deprecated": false
                },
              },
            ],
            "lastActionDate": "15/01/2020",
            "closedOnDate": "12/01/2020",
            "cachedDebtTotal": 12201
          },
          {
            "id": 64,
            "caseRecNumber": "42687993",
            "firstname": "Dolly",
            "surname": "Parton",
            "supervisionCaseOwner": {
              "id": 40,
              "teams": [],
              "displayName": "Supervision Closed Cases"
            },
            "cases": [
              {
                "id": 101,
                "caseRecNumber": "42687993",
                "orderStatus": {
                  "handle": "OPEN",
                  "label": "Open",
                  "deprecated": false
                },
              },
              {
                "id": 102,
                "caseRecNumber": "42687993",
                "orderStatus": {
                  "handle": "ACTIVE",
                  "label": "Active",
                  "deprecated": false
                },
              },
            ],
            "lastActionDate": "01/01/2021",
            "closedOnDate": "12/12/2020",
            "cachedDebtTotal": 10010
          },
        ]
      }
    ],
    "reassign-tasks": [
      {
        "id": 1,
        "assignee":{
          "displayName":"Lay Team 1 - (Supervision)"
        },
      },
      {
        "id": 2,
        "assignee":{
          "displayName":"Pro Team 1 - (Supervision)"
        },
      },
    ],
    "reassign-clients": [
      {
        "id": 1,
        "reAssignName": "LayTeam1 User4",
      }
    ],
    "reassign-deputies": [
      {
        "id": 1,
        "reAssignName": "ProTeam1 User1",
      }
    ],
  }
}
