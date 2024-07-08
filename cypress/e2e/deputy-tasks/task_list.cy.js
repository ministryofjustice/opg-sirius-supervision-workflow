describe("Task list", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/deputy-tasks?team=27");
  });

  it("has column headers", () => {
    cy.get("#workflow-tasks thead > tr > th:nth-child(2)").should("contain", "Task type");
    cy.get("#workflow-tasks thead > tr > th:nth-child(3)").should("contain", "Deputy");
    cy.get("#workflow-tasks thead > tr > th:nth-child(4").should("contain", "Assigned to");
    cy.get("#workflow-tasks thead > tr > th:nth-child(5)").should("contain", "Due date");
  })

  it("has a message to show the team has no tasks", () => {
    cy.get('.moj-team-banner__container > .govuk-form-group > .govuk-select').select("PA Team 1 - (Supervision)");
    cy.get('.govuk-table__cell').should("contain", "The team has no tasks");
  })

  it("should have a table with the column Task type", () => {
    cy.get(".govuk-table__body > :nth-child(1) > :nth-child(2)").should("contain", "PDR follow up")
  })

  it("should have a table with the column Deputy", () => {
    cy.get(".govuk-table__body > :nth-child(1) > :nth-child(3)")
      .should("contain.text", "Mr Fee-paying Deputy")
      .and("contain.text", "Derby - 123456")
      .within(() => {
        cy.get('.govuk-link').should('have.attr', 'href')
          .then(href => {
            expect(href).to.contain("/supervision/deputies/13/tasks");
          })
      })
  })

  it("should not display deputy's town for PA deputies", () => {
    cy.get(".govuk-table__body > :nth-child(2) > :nth-child(3)")
      .should("contain.text", "Mr PRO Deputy")
      .and("contain.text", "654321")
      .and("not.contain.text", "Derby")
      .within(() => {
        cy.get('.govuk-link').should('have.attr', 'href')
          .then(href => {
            expect(href).to.contain("/supervision/deputies/14/tasks");
          })
      })
  })

  it("should have a table with the column Assigned to", () => {
    cy.get(".govuk-table__body > :nth-child(1) > :nth-child(4)").should("contain", "PROTeam1 User1")
  })

  it("should have a table with the column Due date with overdue label", () => {
    cy.get(".govuk-table__body > :nth-child(1) > :nth-child(5)").should("contain", "11/02/2021")
    cy.get('.govuk-tag--red').first().should('contain', "Overdue")
  })
});
