describe("Task list", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/workflow/1");
  });

  it("has column headers", () => {
    cy.get("#workflow-tasks thead > tr > th:nth-child(2)").should("contain", "Task type");
    cy.get("#workflow-tasks thead > tr > th:nth-child(3)").should("contain", "Client");
    cy.get("#workflow-tasks thead > tr > th:nth-child(4)").should("contain", "Case owner");
    cy.get("#workflow-tasks thead > tr > th:nth-child(5)").should("contain", "Assigned to");
    cy.get("#workflow-tasks thead > tr > th:nth-child(6)").should("contain", "Due date");
  })

  it("has a message to show the team has no tasks", () => {
    cy.get('.moj-team-banner__container > .govuk-form-group > .govuk-select').select("Lay Team 1 - (Supervision)");
    cy.get('.govuk-table__cell').should("contain", "The team has no tasks");
  })

  it("should have a table with the column Task type", () => {
    cy.get(".govuk-table__body > :nth-child(1) > :nth-child(2)").should("contain", "Case work - Complaint review")
  })

  it("should have a table with the column Client", () => {
    cy.get(".govuk-table__body > :nth-child(1) > :nth-child(3)").should("contain", "Lizzo Surname")
  })

  it("should have a table with the column Case owner", () => {
    cy.get(".govuk-table__body > :nth-child(1) > :nth-child(4)").should("contain", "Allocations - (Supervision)")
  })

  it("should have a table with the column Assigned to", () => {
    cy.get(".govuk-table__body > :nth-child(1) > :nth-child(5)").should("contain", "Allocations User3")
  })

  it("should have a table with the column Due date", () => {
    cy.get(".govuk-table__body > :nth-child(1) > :nth-child(6)").should("contain", "01/02/2021")
  })

  it("the client name should link to the correct case", () => {
    cy.get(".govuk-table__body > .govuk-table__row > :nth-child(3) > a").should('have.attr', 'href')
        .then(href => {
          expect(href).to.contains("/supervision/#/clients/3333");
        })
    cy.get(".govuk-table__body > .govuk-table__row > :nth-child(3) > a").should('have.attr', 'href')
        .then(href => {
          expect(href).to.contains("/supervision/#/clients/3333");
        })
  })

  it("should display deputy name for PA teams", () => {
    cy.get('.moj-team-banner__container > .govuk-form-group > .govuk-select').select('PA Team 1 - (Supervision)')
    cy.get("#workflow-tasks thead > tr > th:nth-child(4)").should("contain", "Deputy")
    cy.get("#workflow-tasks tbody > tr > td:nth-child(4)").should("contain", "Mr Fee-paying Deputy")
    cy.get("#workflow-tasks tbody > tr > td:nth-child(4) > a").should('have.attr', 'href')
        .then(href => {
          expect(href).to.contains("/supervision/deputies/12");
        })
  })
});
