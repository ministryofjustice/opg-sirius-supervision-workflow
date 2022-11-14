describe("Team Selection", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/workflow/1");
  });

    // it("pulls through my team on the change view bar", () => {
    //   cy.get(".moj-team-banner__container > .govuk-form-group > .govuk-select").should('contain', "Lay Team 1 - (Supervision)")
    // })
    //
    // it("should show the persons team thats logged in", () => {
    //   cy.get(".moj-team-banner__container").should("contain", "Lay Team 1 - (Supervision)")
    // })

    it("can be changed to another team", () => {
        cy.get('.govuk-table__body > :nth-child(1) > :nth-child(5)').should('contain', 'Allocations User3');
        cy.get('.moj-team-banner__container > .govuk-form-group > .govuk-select').select('Lay Team 2 - (Supervision)')
        cy.get(".moj-team-banner__container > .govuk-form-group > .govuk-select").should('contain', "Lay Team 2 - (Supervision)")
        cy.get(".moj-team-banner__container").should("contain", "Lay Team 2 - (Supervision)")
        cy.get('.govuk-table__body > :nth-child(1) > :nth-child(5)').should('contain', 'LayTeam2 User4');
    })

});