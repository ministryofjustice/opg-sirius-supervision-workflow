describe("Reassign Tasks", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.intercept('GET', 'api/v1/teams/*', {
            body: {
                "members": [
                    {
                        "id": 76,
                        "displayName": "LayTeam1 User4",
                    },
                    {
                        "id": 75,
                        "displayName": "LayTeam1 User3",
                    },
                    {
                        "id": 74,
                        "displayName": "LayTeam1 User2",
                    },
                    {
                        "id": 73,
                        "displayName": "LayTeam1 User1",
                    }
                ]
            }})
        cy.visit("/supervision/workflow/1");
    });
    //
    // it("shows me a table of tasks", () => {
    //   cy.get(".govuk-table__body > :nth-child(1) > :nth-child(2) > .govuk-label").contains('Case work - General')
    //   cy.get(":nth-child(1) > :nth-child(3) > .govuk-label").contains('Client Alexander Zacchaeus')
    //   cy.get(":nth-child(1) > :nth-child(4) > .govuk-label").contains('Lay Team 1 - (Supervision)')
    //   cy.get(":nth-child(1) > :nth-child(5) > .govuk-label").contains('LayTeam1 User3')
    // });
    // PATCH /users/27/tasks/0

    it("allows you to assign a task to a team", () => {
        cy.setCookie("success-route", "assignTasksToCasemanager");
        cy.get(":nth-child(1) > :nth-child(1) > .govuk-checkboxes > .govuk-checkboxes__item > #select-task-0").click()
        cy.get("#manage-task").should('be.visible').click()
        cy.get('.moj-manage-tasks__edit-panel > :nth-child(2)').should('be.visible').click()
        cy.get('#assignTeam').select('Pro Team 1 - (Supervision)')

        cy.get('#edit-save').click()

        // cy.intercept('PATCH', 'api/v1/users/*', {statusCode: 204})
        cy.url().should('contain', 'supervision/workflow')
        cy.get("#success-banner").should('be.visible')
        cy.get("#success-banner").contains('1 tasks have been reassigned')
    });
    ///home/runner/work/opg-sirius-supervision-workflow/opg-sirius-supervision-workflow/cypress/screenshots

    // it("allows you to assign multiple tasks to an individual in a team", () => {
    //     cy.setCookie("success-route", "assign-tasks-to-casemanager");
    //     cy.get(":nth-child(1) > :nth-child(1) > .govuk-checkboxes > .govuk-checkboxes__item > #select-task-0").click()
    //     cy.get(":nth-child(2) > :nth-child(1) > .govuk-checkboxes > .govuk-checkboxes__item > #select-task-0").click()
    //     cy.get(":nth-child(5) > :nth-child(1) > .govuk-checkboxes > .govuk-checkboxes__item > #select-task-0").click()
    //     cy.get("#manage-task").should('be.visible').click()
    //     cy.get('.moj-manage-tasks__edit-panel > :nth-child(2)').should('be.visible').click()
    //     cy.get('#assignTeam').select('Pro Team 1 - (Supervision)');
    //     cy.get('#assignCM').select('LayTeam1 User4');
    //     cy.intercept('PATCH', 'api/v1/users/*', {statusCode: 201})
    //     cy.get('#edit-save').click()
    //     cy.intercept('PATCH', 'api/v1/users/*', {statusCode: 201})
    //     cy.url().should('include', 'supervision/workflow')
    //     cy.get("#success-banner").should('be.visible')
    //     cy.get("#success-banner").contains('3 tasks have been reassigned')
    //
    // });
    //
    // it("throws error when task is not assigned to a team", () => {
    //     cy.get(":nth-child(1) > :nth-child(1) > .govuk-checkboxes > .govuk-checkboxes__item > #select-task-0").check('0')
    //     cy.get("#manage-task").click()
    //     cy.get("#edit-save").click()
    //     cy.get(".govuk-error-summary").contains("Please select a team")
    //     cy.wait(5000)
    //     cy.get(".govuk-error-summary").should('not.be.visible')
    // });
    //
    // it("can cancel out of reassigning a task", () => {
    //     cy.get(":nth-child(1) > :nth-child(1) > .govuk-checkboxes > .govuk-checkboxes__item > #select-task-0").check('0')
    //     cy.get("#manage-task").click()
    //     cy.get("#edit-cancel").click()
    //     cy.get(".moj-manage-tasks__edit-panel").should('not.be.visible')
    // });
    //
    // it("does not show the edit panel or manage task by default", () => {
    //     cy.get(".moj-manage-tasks__edit-panel").should('not.be.visible')
    //     cy.get("#manage-task").should('not.be.visible')
    // });

});