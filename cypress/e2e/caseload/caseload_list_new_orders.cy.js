// describe("Caseload list", () => {
//     before(() => {
//         cy.setCookie("Other", "other");
//         cy.setCookie("XSRF-TOKEN", "abcde");
//         cy.visit("/caseload?team=28");
//     });
//
//     it("has column headers", () => {
//         cy.get('[data-cy="Client"]').should("contain", "Client");
//         cy.get('[data-cy="Order made date"]').should("contain", "Order made date");
//         cy.get('[data-cy="Introductory target date"]').should("contain", "Introductory target date");
//         cy.get('[data-cy="Appointment type"]').should("contain", "Appointment type");
//     })
//
//     it("should have a table with the column Client", () => {
//         cy.get('.govuk-table__body > .govuk-table__row > :nth-child(2)').should("contain", "Ro Bot")
//     })
//
//     it("should have a table with the column Order date", () => {
//         cy.get(".govuk-table__body > :nth-child(1) > :nth-child(3)").should("contain", "01/01/2020")
//     })
//
//     it("should have a table with the column Introductory target date", () => {
//         cy.get(".govuk-table__body > :nth-child(1) > :nth-child(4)").should("contain", "21/02/2020")
//     })
//
//     it("should have a table with the column Appointment type", () => {
//         cy.get(".govuk-table__body > :nth-child(1) > :nth-child(5)").should("contain", "Sole")
//     })
// });
