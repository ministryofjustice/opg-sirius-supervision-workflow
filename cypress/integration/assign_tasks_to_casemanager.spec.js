describe("Reassign Tasks", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/");
    });
  
  // it("sends the task to be reassign to someone else", () => {
  // cy.get(":nth-child(1) > :nth-child(1) > .govuk-checkboxes > .govuk-checkboxes__item > #select-task-0").check('0')
  // cy.get("#manage-task").click()
  // cy.get("#assignCM").select('LayTeam1 User11')
  // cy.get("#edit-save").click()
  // cy.get(".moj-banner").contains("1 tasks have been reassigned")
  // cy.wait(5000)
  // cy.get(".moj-banner").should('not.be.visible') 
  // })

  it("can cancel out of reassigning a task", () => {
    cy.get(":nth-child(1) > :nth-child(1) > .govuk-checkboxes > .govuk-checkboxes__item > #select-task-0").check('0')
    cy.get("#manage-task").click()
    cy.get("#edit-cancel").click()
    cy.get(".moj-manage-tasks__edit-panel").should('not.be.visible') 
    })

  it("does not show the edit panel or manage task by default", () => {
    cy.get(".moj-manage-tasks__edit-panel").should('not.be.visible') 
    cy.get("#manage-task").should('not.be.visible') 
    })
  
  });