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

  it("shows user that is logged in within banner", () => {
    cy.contains(".moj-header__link", "case manager");
  });

  const expected = [
    "Power of Attorney",
    "Supervision",
    "Admin",
    "Sign out",
  ];

  it("has working nav links within banner", () => {
    cy.get(".moj-header__navigation-list")
    .children()
    .each(($el, index) => {
        cy.wrap($el).should("contain", expected[index]);
    });
  })

  it("the nav link should contain lpa", () => {
    cy.get(".moj-header__navigation-list > :nth-child(1) > a").should("contain.value", "/lpa")
  })

  it("the nav link should contain supervision", () => {
    cy.get(".moj-header__navigation-list > :nth-child(1) > a").should("contain.value", "/supervision")
  })  
  
  it("the nav link should contain lpa", () => {
    cy.get(".moj-header__navigation-list > :nth-child(2) > a").should("contain.value", "/admin")
  })
  
  it("the nav link should contain logout", () => {
    cy.get(".moj-header__navigation-list > :nth-child(3) > a").should("contain.value", "/auth/logout")
  })
});
