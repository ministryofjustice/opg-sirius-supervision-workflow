describe("Caseload list", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/caseload?team=29");
    });

    it("has column headers", () => {
        cy.get('[data-cy="Client"]').should("contain", "Client");
        cy.get('[data-cy="Deputy type"]').should("contain", "Deputy type");
        cy.get('[data-cy="Case type"]').should("contain", "Case type");
        cy.get('[data-cy="Case owner"]').should("contain", "Case owner");
        cy.get('[data-cy="Status"]').should("contain", "Status");
    })

    it("should have a table with the column Client", () => {
        cy.get('.govuk-table__body > .govuk-table__row > :nth-child(2)').should("contain", "Ro Bot")
    })

    it("should have a table with the column Deputy type", () => {
        cy.get(".govuk-table__body > :nth-child(1) > :nth-child(3)").should("contain", "PA")
    })

    it("should have a table with the column Case type", () => {
        cy.get(".govuk-table__body > :nth-child(1) > :nth-child(4)").should("contain", "Hybrid")
    })

    it("should have a table with the column Case owner", () => {
        cy.get(".govuk-table__body > :nth-child(1) > :nth-child(5)").should("contain", "Not assigned")
    })

    it("should have a table with the column Status", () => {
        cy.get(".govuk-table__body > :nth-child(1) > :nth-child(6)").should("contain", "Active")
    })
});
