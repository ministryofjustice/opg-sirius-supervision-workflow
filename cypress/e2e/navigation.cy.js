describe("Navigation", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/workflow/1");
    });

    it("supervision link is clicked", () => {
        cy.get(':nth-child(1) > .moj-header__navigation-link').click()
        cy.url().should('include', 'localhost:8080/supervision')
    })

    it("lpa link is clicked", () => {
        cy.get(':nth-child(2) > .moj-header__navigation-link').click()
        cy.url().should('include', 'localhost:8080/lpa')
    })

    it("log out link is clicked", () => {
        cy.get(':nth-child(3) > .moj-header__navigation-link').click()
        cy.url().should('include', 'localhost:8080/auth/logout')
    })
});