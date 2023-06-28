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
    cy.get('.moj-header__logo > .moj-header__link').should('contain.text', 'OPG');
    cy.contains(".moj-header__link", "Sirius");
  });
});
