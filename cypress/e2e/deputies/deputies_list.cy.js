describe("Deputies list", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/deputies?team=27");
    });

    it("has column headers", () => {
        cy.get("th:nth-child(2)").should("contain", "Deputy");
        cy.get("th:nth-child(3)").should("contain", "Executive case manager");
        cy.get("th:nth-child(4)").should("contain", "Active clients");
        cy.get("th:nth-child(5)").should("contain", "Non-compliance");
        cy.get("th:nth-child(6)").should("contain", "Assurance visits");
    })

    it("has column values", () => {
        cy.get(".govuk-table__body .govuk-table__row:nth-child(1)").within(() => {
            cy.get("td:nth-child(2)").should("contain.text", "Mr Fee-paying Deputy").and("contain.text", "Derby - 123456")
            cy.get("td:nth-child(3)").should("contain.text", "PROTeam1 User1")
            cy.get("td:nth-child(4)").should("contain.text", "100")
            cy.get("td:nth-child(5)").should("contain.text", "10 (10%)")
            cy.get("td:nth-child(6)").should("contain.text", "26/05/2023").and("contain.text", "Low risk")
        })
    })

    it("can be sorted by Non-compliance and Deputy", () => {
        cy.url().should("not.contain", "order-by").and("not.contain", "sort")

        cy.get("th:nth-child(5) button").click()
        cy.get("th:nth-child(5)").should("have.attr", "aria-sort", "ascending")
        cy.url().should("contain", "order-by=noncompliance&sort=asc")

        cy.get("th:nth-child(5) button").click()
        cy.get("th:nth-child(5)").should("have.attr", "aria-sort", "descending")
        cy.url().should("contain", "order-by=noncompliance&sort=desc")

        cy.get("th:nth-child(2) button").click()
        cy.get("th:nth-child(2)").should("have.attr", "aria-sort", "ascending")
        cy.url().should("contain", "order-by=deputy&sort=asc")

        cy.get("th:nth-child(2) button").click()
        cy.get("th:nth-child(2)").should("have.attr", "aria-sort", "descending")
        cy.url().should("contain", "order-by=deputy&sort=desc")
    })
});
