describe("Work flow", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/workflow");
    });

    it("shows user that is logged in within banner", () => {
      cy.contains(".moj-header__link", "system admin");
    });


    const expected = [
      "Supervision",
      "LPA",
      "Log out",
  ];

    it("has working nav links within banner", () => {
      cy.get(".moj-header__navigation-list").each(($el, index) => {
        cy.wrap($el).within(() => {
          cy.get(".moj-header__navigation-link").first().should(
                "have.text",
                expected[index]
          )
        })
      })
    })
   
});