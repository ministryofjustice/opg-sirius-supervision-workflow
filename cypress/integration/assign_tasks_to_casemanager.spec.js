describe("Reassign Tasks", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/");
    });
  
   it("sends the task to be reassign to someone else", () => {
    cy.get(":nth-child(1) > :nth-child(1) > .govuk-checkboxes > .govuk-checkboxes__item > #select-task-0").check('0')
    cy.get("#edit-task").click()
    cy.get("#assignCM").select('LayTeam1 User11')
    cy.get("#edit-save").click()
    cy.get(".moj-banner").contains("1 tasks have been reassigned")
    
  })
  
  });