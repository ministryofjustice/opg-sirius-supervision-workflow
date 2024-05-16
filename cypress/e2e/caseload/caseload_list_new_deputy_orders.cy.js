describe("Caseload list", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    it("should redirect me to client tasks as there is no lay new deputy orders caseload page ", () => {
        cy.visit("/caseload?team=29");
        cy.get('.moj-team-banner__container > .govuk-form-group > .govuk-select').select('Lay Team - New Deputy Orders')
        cy.url().should('contain', '/client-tasks?team=28')
    })
});
