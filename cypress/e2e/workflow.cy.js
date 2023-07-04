describe("Workflow", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/client-tasks");
  });

  it("homepage redirects to client tasks page", () => {
    cy.visit("/caseload?team=21");
    cy.url().should('not.contain', '/client-tasks')
    cy.visit("/");
    cy.url().should('contain', '/client-tasks')
  });

  it("should load header template within banner", () => {
      cy.get('.govuk-header__link--homepage').should('contain.text', 'OPG');
      cy.get('.govuk-header__service-name').should('contain.text', 'Sirius')
  });
});
