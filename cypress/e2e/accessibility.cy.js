import "cypress-axe";

describe("Accessibility test client tasks", () => {
    before(() => {
        cy.visit('/');
        cy.url().should('contain', '/client-tasks')
        cy.injectAxe();
    });

    it("Should have no accessibility violations",() => {
        cy.debug();
        cy.checkA11y();
    });
});

describe("Accessibility test caseload", () => {
    before(() => {
        cy.visit('/caseload?team=21');
        cy.injectAxe();
    });

    it("Should have no accessibility violations",() => {
        cy.checkA11y();
    });
});