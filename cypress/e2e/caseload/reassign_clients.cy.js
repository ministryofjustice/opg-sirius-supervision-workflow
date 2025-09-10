describe("Reassign clients", () => {
  const pages = [
    "/caseload?team=21",
    "/caseload?team=28"
  ]

  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
      cy.intercept('supervision-api/v1/teams/21', {
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

it("allows you to reassign a client for normal lay team", () => {
    cy.visit("/caseload?team=21")
    cy.setCookie("success-route", "/reassign-clients/1");
    cy.get('#select-client-63').click();
    cy.get("#manage-client").should('be.visible').click();
    cy.get('.moj-manage-list__edit-panel > :nth-child(2)').should('be.visible').click()
    cy.get('#assignTeam').select('Lay Team 1 - (Supervision)');
    cy.intercept('PATCH', 'supervision-api/v1/users/*', {statusCode: 204})
    cy.get('#assignCM option:contains(LayTeam1 User3)').should('not.exist')
    cy.get('#assignCM option:contains(LayTeam1 User1)').should('exist')
    cy.get('#assignCM').select('LayTeam1 User1');
    cy.get('#edit-save').click()
    cy.getCookies()
      .should('have.length', 4)
      .then((cookies) => {
        expect(cookies[0]).to.have.property('name', 'successMessageStore'),
        expect(cookies[1]).to.have.property('name', 'Other'),
        expect(cookies[2]).to.have.property('name', 'XSRF-TOKEN'),
        expect(cookies[3]).to.have.property('name', 'success-route')
      })
    cy.get("#success-banner").should('exist')
    cy.get("#success-banner").should('be.visible')
    cy.get("#success-banner").contains('You have reassigned ')
  });

  it("allows you to reassign a client for lay new deputy team", () => {
      cy.visit("/caseload?team=28")
      cy.setCookie("success-route", "/reassign-clients/1");
      cy.get('#select-client-63').click();
      cy.get("#manage-client").should('be.visible').click();
      cy.get('.moj-manage-list__edit-panel > :nth-child(2)').should('be.visible').click()
      cy.get('#assignTeam').select('Lay Team 1 - (Supervision)');
      cy.intercept('PATCH', 'supervision-api/v1/users/*', {statusCode: 204})
      cy.get('#assignCM option:contains(LayTeam1 User3)').should('not.exist')
      cy.get('#assignCM option:contains(LayTeam1 User4)').should('exist')
      cy.get('#assignCM').select('LayTeam1 User4');
      cy.get('#edit-save').click()
      cy.getCookies()
        .should('have.length', 4)
        .then((cookies) => {
          expect(cookies[0]).to.have.property('name', 'successMessageStore'),
          expect(cookies[0]).to.have.property('value', 'MTc1NzUxNzMyMHxEWDhFQVFMX2dBQUJFQUVRQUFBRV80QUFBQT09fHOuzLICRji8VovEjFjIcXm202Fm1JSUV6Rv3RoBB4wZ'),
          expect(cookies[1]).to.have.property('name', 'Other'),
          expect(cookies[2]).to.have.property('name', 'XSRF-TOKEN'),
          expect(cookies[3]).to.have.property('name', 'success-route')
        })
      cy.get("#success-banner").should('exist')
      cy.get("#success-banner").should('be.visible')
      cy.get("#success-banner").contains('You have reassigned ')
  });
})