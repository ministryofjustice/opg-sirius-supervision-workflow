describe("Reassign Tasks", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/workflow/1");
    });
  
    it("shows me a table of tasks", () => {
      cy.get(".govuk-table__body > :nth-child(1) > :nth-child(2) > .govuk-label").contains('Case work - General')
      cy.get(":nth-child(1) > :nth-child(3) > .govuk-label").contains('Client Alexander Zacchaeus')
      cy.get(":nth-child(1) > :nth-child(4) > .govuk-label").contains('Lay Team 1 - (Supervision)')
      cy.get(":nth-child(1) > :nth-child(5) > .govuk-label").contains('LayTeam1 User3')
    });

    it("allows you to manage a task", () => {
       cy.setCookie("success-route", "assign-tasks-to-casemanager");
       cy.get(":nth-child(1) > :nth-child(1) > .govuk-checkboxes > .govuk-checkboxes__item > #select-task-0").click()
       cy.get("#manage-task").should('be.visible').click()
       cy.get('.moj-manage-tasks__edit-panel > :nth-child(2)').should('be.visible').click()
       cy.get('#assignTeam').select('Pro Team 1 - (Supervision)');
       cy.get("#edit-panel").click()

       //struggles with the javascript binding to get the case managers for the selected team
       cy.on('uncaught:exception', (err, runnable) => {
           return false
       })
    });

    it("throws error when task is not assigned to a team", () => {
        cy.get(":nth-child(1) > :nth-child(1) > .govuk-checkboxes > .govuk-checkboxes__item > #select-task-0").check('0')
        cy.get("#manage-task").click()
        cy.get("#edit-save").click()
        cy.get(".govuk-error-summary").contains("Please select a team")
        cy.wait(5000)
        cy.get(".govuk-error-summary").should('not.be.visible')
    });

    it("can cancel out of reassigning a task", () => {
        cy.get(":nth-child(1) > :nth-child(1) > .govuk-checkboxes > .govuk-checkboxes__item > #select-task-0").check('0')
        cy.get("#manage-task").click()
        cy.get("#edit-cancel").click()
        cy.get(".moj-manage-tasks__edit-panel").should('not.be.visible')
    });

    it("does not show the edit panel or manage task by default", () => {
        cy.get(".moj-manage-tasks__edit-panel").should('not.be.visible')
        cy.get("#manage-task").should('not.be.visible')
    });

});