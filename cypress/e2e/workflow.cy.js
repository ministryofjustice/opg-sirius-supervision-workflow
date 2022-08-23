describe("Work flow", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/");
  });

  it("shows user that is logged in within banner", () => {
    cy.contains(".moj-header__link", "case manager");
  });

  const expected = [
    "Supervision",
    "LPA",
    "Log out",
];

  it("has working nav links within banner", () => {
    cy.get(".moj-header__navigation-list")
    .children()
    .each(($el, index) => {
        cy.wrap($el).should("contain", expected[index]);
    });
  })

  it("the nav link should contain supervision", () => {
    cy.get(".moj-header__navigation-list > :nth-child(1) > a").should("have.attr", "href", "http://localhost:8080/supervision")
  })  
  
  it("the nav link should contain lpa", () => {
    cy.get(".moj-header__navigation-list > :nth-child(2) > a").should("have.attr", "href", "http://localhost:8080/lpa")
  })
  
  it("the nav link should contain logout", () => {
    cy.get(".moj-header__navigation-list > :nth-child(3) > a").should("have.attr", "href", "http://localhost:8080/auth/logout")
  }) 
});