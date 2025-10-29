describe("Bonds page redirect", () => {
    before(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    })

    it("Non-allocations teams are redirected to the Client tasks list", () => {
        cy.visit("/bonds?team=21");
        cy.url().should('contain', '/client-tasks?team=21')
    })

    it("Switching to another team from the Bonds page should redirect to the Client tasks list", () => {
        cy.visit("/bonds?team=13");
        cy.get('.moj-team-banner__container > .govuk-form-group > .govuk-select').select('Lay Team 1 - (Supervision)')
        cy.url().should('contain', '/client-tasks?team=21')
    })
})
