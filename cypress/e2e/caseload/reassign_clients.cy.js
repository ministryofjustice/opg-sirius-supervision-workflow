describe("Reassign clients", () => {
  const pages = [
    "/caseload?team=21",
    "/caseload?team=28"
  ]

  before(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
  });

  it("can cancel out of reassigning a client", () => {
    pages.forEach((page) => {
      cy.visit(page)
      cy.get('#select-client-63').click();
      cy.get("#manage-client").should('be.visible').click();
      cy.get("#edit-cancel").click();
      cy.get(".moj-manage-list__edit-panel").should('not.be.visible');
    })
  });

  it("allows you to reassign a client", () => {
    pages.forEach((page) => {
      cy.visit(page)
      cy.intercept('api/v1/teams/21', {
        body: {
          "members": [
            {
              "id": 76,
              "displayName": "LayTeam1 User4",
              "suspended": false,
            },
            {
              "id": 75,
              "displayName": "LayTeam1 User3",
              "suspended": true,
            },
            {
              "id": 74,
              "displayName": "LayTeam1 User2",
              "suspended": false,
            },
            {
              "id": 73,
              "displayName": "LayTeam1 User1",
              "suspended": false,
            }
          ]
        }
      })
      cy.setCookie("success-route", "/reassign-clients/1");
      cy.get('#select-client-63').click();
      cy.get("#manage-client").should('be.visible').click();
      cy.get('.moj-manage-list__edit-panel > :nth-child(2)').should('be.visible').click()
      cy.get('#assignTeam').select('Lay Team 1 - (Supervision)');
      cy.intercept('PATCH', 'api/v1/users/*', {statusCode: 204})
      cy.get('#assignCM option:contains(LayTeam1 User3)').should('not.exist')
      cy.get('#assignCM option:contains(LayTeam1 User4)').should('exist')
      cy.get('#assignCM').select('LayTeam1 User4');
      cy.get('#edit-save').click()
      cy.get("#success-banner").should('be.visible')
      cy.get("#success-banner").contains('You have reassigned ')
    })
  });
})