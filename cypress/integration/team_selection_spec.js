describe("Team Selection", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/workflow");
  });

 it("pulls through my team on the change view bar", () => {
  cy.get("#change-team").should('contain', "Lay Team 1 - (Supervision)")
})

});