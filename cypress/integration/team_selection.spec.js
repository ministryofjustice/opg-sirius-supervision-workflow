describe("Team Selection", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/");
  });

 it("pulls through my team on the change view bar", () => {
  cy.get("#team-banner-container > .govuk-form-group > .govuk-select").should('contain', "Go TaskForce")
})

it("should show the persons team thats logged in", () => {
  cy.get("#hook-team-name").should("contain", "Go TaskForce")
})

});