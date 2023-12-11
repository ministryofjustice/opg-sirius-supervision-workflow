describe("Caseload list", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/caseload?team=40");
    });

    it("has column headers", () => {
        cy.get('[data-cy="Client"]').should("contain", "Client");
        cy.get('[data-cy="Closed-on-date"]').should("contain", "Closed on date");
        cy.get('[data-cy="Last-action-date"]').should("contain", "Last action date");
        cy.get('[data-cy="Debt"]').should("contain", "Debt");
        cy.get('[data-cy="Status"]').should("contain", "Status");
    })

    it("should have a table with the column Client", () => {
        cy.get('.govuk-table__body > .govuk-table__row > :nth-child(2)').should("contain", "Ro Bot")
    })

    it("should have a table with the column Closed On Date which shows most recently closed orders date", () => {
        cy.get(".govuk-table__body > :nth-child(1) > :nth-child(3)").should("contain", "12/01/2020")
    })

    it("should have a table with the column Last Action date", () => {
        cy.get(".govuk-table__body > :nth-child(1) > :nth-child(4)").should("contain", "15/01/2020")
    })

    it("should have a table with the column Debt", () => {
        cy.get(".govuk-table__body > :nth-child(1) > :nth-child(5)").should("contain", "£122.01")
    })

    it("should have a table with the column Status", () => {
        cy.get(".govuk-table__body > :nth-child(1) > :nth-child(6)").should("contain", "Closed")
    })
});
