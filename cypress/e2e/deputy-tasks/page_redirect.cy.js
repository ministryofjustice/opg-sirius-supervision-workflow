describe("Caseload page redirect", () => {
    before(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    })

    it("Lay teams are redirected to the Client tasks list", () => {
        cy.visit("/deputy-tasks?team=21");
        cy.url().should('contain', '/client-tasks?team=21')
    })

    it("Switching to a Lay team from the Deputy tasks page should redirect to the Client tasks list", () => {
        cy.visit("/deputy-tasks?team=27");
        cy.get('.moj-team-banner__container > .govuk-form-group > .govuk-select').select('Lay Team 1 - (Supervision)')
        cy.url().should('contain', '/client-tasks?team=21')
    })

    it("Switching to another Pro/PA team from the Deputy tasks page should not redirect to the Client tasks list", () => {
        cy.visit("/deputy-tasks?team=27");
        cy.get('.moj-team-banner__container > .govuk-form-group > .govuk-select').select('PA Team 1 - (Supervision)')
        cy.url().should('contain', '/deputy-tasks?team=24')
    })
})
