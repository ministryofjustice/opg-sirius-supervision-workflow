// import taskListFixture from 'cypress/fixtures/task_list_fake_data.json'

// describe("Task list fake data", () => {
//     beforeEach(() => {
//       cy.setCookie("Other", "other");
//       cy.setCookie("XSRF-TOKEN", "abcde");
//       cy.visit("/supervision/workflow");
//     });

//     it("has column headers", () => {
//         cy.contains("Task type");
//         cy.contains("Client");
//         cy.contains("Case owner");
//         cy.contains("Assigned to");
//         cy.contains("Due date");
//         cy.contains("Actions");
//       })
    
//       const expected = [
//         "",
//         "Case work - General",
//         "Client Alexander Zacchaeus Client Wolfeschlegelsteinhausenbergerdorff caseRecNumber",
//         "Assignee Duke Clive Henry Hetley Junior Jones",
//         "Assignee Duke Clive Henry Hetley Junior Jones Supervision - Team - Name",
//         "01/02/2021",
//         "Open case",
//     ];
    
//     it("should have data in the table", () => {
//       cy.get(".govuk-table__body > .govuk-table__row")
//         .children()
//         .each(($el, index) => {
//             cy.wrap($el).should("contain", expected[index]);
//         });
//       })
      
//       it("the button should have a link to the correct case", () => {
//         cy.get(".govuk-table__body > .govuk-table__row > :nth-child(7) > a").should('have.attr', 'href', 'http://localhost:8080/supervision/#/clients/3333')
//       })
//   });


  const urls = ['https://localhost:8080', 'https://localhost:8888']
    describe('Logo', () => {
    urls.forEach((url) => {
        it(`Should display logo on ${url}`, () => {
            if(Cypress._.isArray(url) == 'https://localhost:8080') {
                cy.get(':nth-child(1) > label > input').type("lay1-1@opgtest.com")
                cy.get(':nth-child(2) > label > input').type("Password1")
                cy.get('#submitbutton').click()
                cy.visit("https://localhost:8888");
                // cy.fixture("cypress/fixtures/task_list_fake_data.json")
            }
        it("the button should have a link to the correct case", () => {
        cy.get(".govuk-table__body > .govuk-table__row > :nth-child(7) > a").should('have.attr', 'href', 'http://localhost:8080/supervision/#/clients/3333')
      })
    })
    })
})
  