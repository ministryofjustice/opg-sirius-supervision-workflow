describe("Work flow", () => {
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

  const expected = [
    "",
    "Case work - General",
    "Client Alexander Zacchaeus Client Wolfeschlegelsteinhausenbergerdorff caseRecNumber",
    "Assignee Duke Clive Henry Hetley Junior Jones",
    "Assignee Duke Clive Henry Hetley Junior Jones Supervision - Team - Name",
    "01/02/2021",
    "Open case",
];

it("should have data in the table", () => {
  cy.get(".govuk-table__body > .govuk-table__row")
    .children()
    .each(($el, index) => {
        cy.wrap($el).should("contain", expected[index]);
    });
  })
  
  it("the button should have a link to the correct case", () => {
    cy.get(".govuk-table__body > .govuk-table__row > :nth-child(7) > a").should('have.attr', 'href', 'http://localhost:8080/supervision/#/clients/3333')
  })
});
