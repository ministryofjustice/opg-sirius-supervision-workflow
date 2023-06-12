describe("Caseload page redirect", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    })

    it("Non-Lay teams are redirected to the Client tasks list", () => {
        cy.visit("/caseload?team=13");
        cy.url().should('contain', '/client-tasks?team=13')
    })

    it("Switching to a non-Lay team from the Caseload page should redirect to the Client tasks list", () => {
        cy.visit("/caseload?team=21");
        cy.get('.moj-team-banner__container > .govuk-form-group > .govuk-select').select('Pro Team 1 - (Supervision)')
        cy.url().should('contain', '/client-tasks?team=27')
    })

    it("Switching to another Lay team from the Caseload page should not redirect to the Client tasks list", () => {
        cy.visit("/caseload?team=21");
        cy.get('.moj-team-banner__container > .govuk-form-group > .govuk-select').select('Lay Team 2 - (Supervision)')
        cy.url().should('contain', '/caseload?team=22')
    })
})
