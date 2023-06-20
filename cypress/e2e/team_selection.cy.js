describe("Team Selection", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/client-tasks");
  });

    it("can be changed to another team", () => {
        cy.get('.govuk-table__body > :nth-child(1) > :nth-child(5)').should('contain', 'Allocations User3');
        cy.get('.moj-team-banner__container > .govuk-form-group > .govuk-select').select('Lay Team 2 - (Supervision)')
        cy.get(".moj-team-banner__container > .govuk-form-group > .govuk-select").should('contain', "Lay Team 2 - (Supervision)")
        cy.get(".moj-team-banner__container").should("contain", "Lay Team 2 - (Supervision)")
        cy.get('.govuk-table__body > :nth-child(1) > :nth-child(5)').should('contain', 'LayTeam2 User4');
    })

    it("contains options for combined Lay and Pro teams", () => {
        cy.get(".moj-team-banner__container > .govuk-form-group > .govuk-select").should('contain', "Lay deputy team")
        cy.get(".moj-team-banner__container > .govuk-form-group > .govuk-select").should('contain', "Professional deputy team")
    })
});
