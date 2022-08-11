describe("Team Selection", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.intercept('api/v1/teams/*', {
          body: {
              "members": [
                  {
                      "id": 76,
                      "displayName": "LayTeam1 User4",
                  },
                  {
                      "id": 75,
                      "displayName": "LayTeam1 User3",
                  },
                  {
                      "id": 74,
                      "displayName": "LayTeam1 User2",
                  },
                  {
                      "id": 73,
                      "displayName": "LayTeam1 User1",
                  }
              ]
          }})
      cy.visit("/supervision/workflow/1");
  });

//  it("pulls through my team on the change view bar", () => {
//   cy.get(".moj-team-banner__container > .govuk-form-group > .govuk-select").should('contain', "Lay Team 1 - (Supervision)")
// })

// it("should show the persons team thats logged in", () => {
//   cy.get(".moj-team-banner__container").should("contain", "Lay Team 1 - (Supervision)")
// })

});