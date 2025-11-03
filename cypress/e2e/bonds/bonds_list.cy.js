describe("Bonds list", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/bonds?team=13");
    });

    it("has column headers", () => {
        cy.get("th:nth-child(1)").should("contain", "Order number");
        cy.get("th:nth-child(2)").should("contain", "Client name");
        cy.get("th:nth-child(3)").should("contain", "Bond company");
        cy.get("th:nth-child(4)").should("contain", "Bond amount");
        cy.get("th:nth-child(5)").should("contain", "Bond reference");
        cy.get("th:nth-child(6)").should("contain", "Date issued");
    })

    it("has column values", () => {
        cy.get(".govuk-table__body .govuk-table__row:nth-child(1)").within(() => {
            cy.get("td:nth-child(1)").should("contain.text", "12345678")
            cy.get("td:nth-child(2)").should("contain.text", "John Smith")
            cy.get("td:nth-child(3)").should("contain.text", "Marsh")
            cy.get("td:nth-child(4)").should("contain.text", "£100.00")
            cy.get("td:nth-child(5)").should("contain.text", "BOND-001")
            cy.get("td:nth-child(6)").should("contain.text", "01/01/2024")
        })
    })
});
