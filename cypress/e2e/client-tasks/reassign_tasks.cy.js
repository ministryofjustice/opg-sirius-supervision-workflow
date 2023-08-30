describe("Reassign Tasks", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.setCookie("success-route", "/reassign-tasks/1");

        cy.intercept('api/v1/teams/21', {
            body: {
                "members": [
                    {
                        "id": 76,
                        "displayName": "LayTeam1 User4",
                        "suspended": false,
                    },
                    {
                        "id": 75,
                        "displayName": "LayTeam1 User3",
                        "suspended": true,
                    },
                    {
                        "id": 74,
                        "displayName": "LayTeam1 User2",
                        "suspended": false,
                    },
                    {
                        "id": 73,
                        "displayName": "LayTeam1 User1",
                        "suspended": false,
                    }
                ]
            }})
        cy.visit("/client-tasks");
    });

    it("allows you to assign a task to a team and retains pagination and filters", () => {
        cy.visit('/client-tasks?testVar=testVal');
        cy.get("#select-task-1").click()
        cy.get("#manage-task").should('be.visible').click()
        cy.get('.moj-manage-list__edit-panel > :nth-child(2)').should('be.visible').click()
        cy.get('#assignTeam').select('Lay Team 1 - (Supervision)')
        cy.intercept('PATCH', 'api/v1/users/*', {statusCode: 204})
        cy.get('#edit-save').click()
        cy.get("#success-banner").should('be.visible')
        cy.get("#success-banner").contains('You have assigned 1 task(s) to Lay Team 1 - (Supervision)')
        cy.url().should('contain', '/client-tasks?testVar=testVal')
    });

    it("allows you to assign multiple tasks to an individual in a team", () => {
        cy.get("#select-task-1").click()
        cy.get("#select-task-2").click()
        cy.get("#select-task-5").click()
        cy.get("#manage-task").should('be.visible').click()
        cy.get('.moj-manage-list__edit-panel > :nth-child(2)').should('be.visible').click()
        cy.get('#assignTeam').select('Lay Team 1 - (Supervision)');
        cy.intercept('PATCH', 'api/v1/users/*', {statusCode: 204})
        cy.get('#assignCM option:contains(LayTeam1 User3)').should('not.exist')
        cy.get('#assignCM option:contains(LayTeam1 User4)').should('exist')
        cy.get('#assignCM').select('LayTeam1 User4');
        cy.get('#edit-save').click()
        cy.get("#success-banner").should('be.visible')
        cy.get("#success-banner").contains('You have assigned 3 task(s) to Lay Team 1 - (Supervision)')
    });

    it("can cancel out of reassigning a task", () => {
        cy.get("#select-task-1").check('1')
        cy.get("#manage-task").click()
        cy.get("#edit-cancel").click()
        cy.get(".moj-manage-list__edit-panel").should('not.be.visible')
    });

    it("Only set the priority for a task", () => {
        cy.visit('/client-tasks?testVar=testVal');
        cy.get("#select-task-1").click()
        cy.get("#manage-task").should('be.visible').click()
        cy.get('#priority').select('Yes')
        cy.get('#edit-save').click()
        cy.get("#success-banner").should('be.visible')
        cy.get("#success-banner").contains('You have assigned 1 task(s) as a priority')
    })

    it("Reassign and set the priority for a task", () => {
        cy.get("#select-task-1").click()
        cy.get("#manage-task").should('be.visible').click()
        cy.get('#assignTeam').select('Lay Team 1 - (Supervision)');
        cy.get('#priority').select('Yes')
        cy.get('#edit-save').click()
        cy.get("#success-banner").should('be.visible')
        cy.get("#success-banner").contains('You have assigned 1 task(s) to Lay Team 1 - (Supervision) as a priority')
    })
});