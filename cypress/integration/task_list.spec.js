describe("Task list", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/workflow");
  });

  it("has column headers", () => {
    cy.contains("Task type");
    cy.contains("Client");
    cy.contains("Case owner");
    cy.contains("Assigned to");
    cy.contains("Due date");
    cy.contains("Actions");
  })

it("should have a table with the column Task type", () => {
  cy.get(".govuk-table__body > :nth-child(1) > :nth-child(2)").should("contain", "Case work - General")
})

it("should have a table with the column Client", () => {
  cy.get(".govuk-table__body > :nth-child(1) > :nth-child(3)").should("contain", "Client Alexander Zacchaeus Client Wolfeschlegelsteinhausenbergerdorff caseRecNumber")
})

it("should have a table with the column Case owner", () => {
  cy.get(".govuk-table__body > :nth-child(1) > :nth-child(4)").should("contain", "Assignee Duke Clive Henry Hetley Junior Jones")
})

it("should have a table with the column Assigned to", () => {
  cy.get(".govuk-table__body > :nth-child(1) > :nth-child(5)").should("contain", "Assignee Duke Clive Henry Hetley Junior Jones Supervision - Team - Name")
})

it("should have a table with the column Due date", () => {
  cy.get(".govuk-table__body > :nth-child(1) > :nth-child(6)").should("contain", "01/02/2021")
})

it("should have a table with the column Actions", () => {
  cy.get(".govuk-table__body > :nth-child(1) > :nth-child(7)").should("contain", "Open case")
})
  
it("the button should have a link to the correct case", () => {
    cy.get(".govuk-table__body > .govuk-table__row > :nth-child(7) > a").should('have.attr', 'href', 'http://localhost:8080/supervision/#/clients/3333')
  })
});
