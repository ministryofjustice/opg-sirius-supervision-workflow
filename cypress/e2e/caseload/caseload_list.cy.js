describe("Caseload list", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    it("has column headers", () => {
        cy.visit("/client-tasks?team=21");
        cy.get(':nth-child(2) > .moj-sub-navigation__link').click()
        cy.get("#caseload-list thead > tr > th:nth-child(2)").should("contain", "Client");
        cy.get("#caseload-list thead > tr > th:nth-child(3)").should("contain", "Report due date");
        cy.get("#caseload-list thead > tr > th:nth-child(4)").should("contain", "Case owner");
        cy.get("#caseload-list thead > tr > th:nth-child(5)").should("contain", "Supervision level");
        cy.get("#caseload-list thead > tr > th:nth-child(6)").should("contain", "Status");

    })

    it("should have a table with the column Client", () => {
        cy.get('.govuk-table__body > .govuk-table__row > :nth-child(2)').should("contain", "Ro Bot")
    })

    it("should have a table with the column Report due date", () => {
        cy.get(".govuk-table__body > :nth-child(1) > :nth-child(3)").should("contain", "21/12/2023")
    })

    it("should have a table with the column Case owner", () => {
        cy.get(".govuk-table__body > :nth-child(1) > :nth-child(4)").should("contain", "Lay Team 1 - (Supervision)")
    })

    it("should have a table with the column Supervision level", () => {
        cy.get(".govuk-table__body > :nth-child(1) > :nth-child(5)").should("contain", "Minimal")
    })

    it("should have a table with the column Status", () => {
        cy.get(".govuk-table__body > :nth-child(1) > :nth-child(6)").should("contain", "Closed")
    })
});
