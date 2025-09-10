//describe("PA Deputies list", () => {
//    beforeEach(() => {
//        cy.setCookie("Other", "other");
//        cy.setCookie("XSRF-TOKEN", "abcde");
//        cy.visit("/deputies?team=24");
//    });
//
//    it("has column headers", () => {
//        cy.get("th:nth-child(2)").should("contain", "Deputy");
//        cy.get("th:nth-child(3)").should("contain", "Executive case manager");
//        cy.get("th:nth-child(4)").should("contain", "Active clients");
//        cy.get("th:nth-child(5)").should("contain", "Non-compliance");
//        cy.get("th:nth-child(6)").should("contain", "Assurance visits");
//    })
//
//    it("has column values", () => {
//        cy.get(".govuk-table__body .govuk-table__row:nth-child(1)").within(() => {
//            cy.get("td:nth-child(2)").should("contain.text", "Mr PA Deputy")
//            cy.get("td:nth-child(3)").should("contain.text", "PA Team 1 - (Supervision)")
//            cy.get("td:nth-child(4)").should("contain.text", "81")
//            cy.get("td:nth-child(5)").should("contain.text", "34 (42%)")
//        })
//    })
//
//    it("can be sorted by Deputy, Active clients, Non-compliance & Assurance", () => {
//        cy.url().should("not.contain", "order-by").and("not.contain", "sort")
//
//        cy.get("a > button").contains('Deputy').click()
//        cy.get("th:nth-child(2)").should("have.attr", "aria-sort", "descending")
//        cy.url().should("contain", "order-by=deputy&sort=desc")
//
//        cy.get("a > button").contains('Deputy').click()
//        cy.get("th:nth-child(2)").should("have.attr", "aria-sort", "ascending")
//        cy.url().should("contain", "order-by=deputy&sort=asc")
//
//        cy.get("a > button").contains('Active clients').click()
//        cy.get("th:nth-child(4)").should("have.attr", "aria-sort", "ascending")
//        cy.url().should("contain", "order-by=activeclients&sort=asc")
//
//        cy.get("a > button").contains('Active clients').click()
//        cy.get("th:nth-child(4)").should("have.attr", "aria-sort", "descending")
//        cy.url().should("contain", "order-by=activeclients&sort=desc")
//
//        cy.get("a > button").contains('Non-compliance').click()
//        cy.get("th:nth-child(5)").should("have.attr", "aria-sort", "ascending")
//        cy.url().should("contain", "order-by=noncompliance&sort=asc")
//
//        cy.get("a > button").contains('Non-compliance').click()
//        cy.get("th:nth-child(5)").should("have.attr", "aria-sort", "descending")
//        cy.url().should("contain", "order-by=noncompliance&sort=desc")
//
//        cy.get("a > button").contains('Assurance visits').click()
//        cy.get("th:nth-child(6)").should("have.attr", "aria-sort", "ascending")
//        cy.url().should("contain", "order-by=assurance&sort=asc")
//
//        cy.get("a > button").contains('Assurance visits').click()
//        cy.get("th:nth-child(6)").should("have.attr", "aria-sort", "descending")
//        cy.url().should("contain", "order-by=assurance&sort=desc")
//    })
//});
//
//describe("Pro Deputies list", () => {
// beforeEach(() => {
//        cy.setCookie("Other", "other");
//        cy.setCookie("XSRF-TOKEN", "abcde");
//        cy.visit("/deputies?team=27");
//    });
//
//    it("has additional column headers", () => {
//        cy.get(".moj-team-banner__container > .govuk-form-group > .govuk-select").select("Professional Deputy Team")
//        cy.get("th:nth-child(3)").should("contain", "Firm");
//
//        cy.get(".govuk-table__body .govuk-table__row:nth-child(1)").within(() => {
//            cy.get("td:nth-child(2)").should("contain.text", "Mr Fee-paying Deputy")
//                .and("contain.text", "Derby - 123456")
//                .and("contain.text", "Panel Deputy")
//           cy.get("td:nth-child(3)").should("contain.text", "Krusty Krabs")
//                .and("contain.text", "789456123")
//            cy.get("td:nth-child(4)").should("contain.text", "PROTeam1 User1")
//            cy.get("td:nth-child(5)").should("contain.text", "100")
//            cy.get("td:nth-child(6)").should("contain.text", "10 (10%)")
//            cy.get("td:nth-child(7)").should("contain.text", "26/05/2023").and("contain.text", "Low risk")
//        })
//    })
//});