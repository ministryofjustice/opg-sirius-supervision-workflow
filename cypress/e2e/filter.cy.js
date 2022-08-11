describe("Reassign Tasks", () => {

  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.window().then((win) => {
      win.sessionStorage.clear()
    })
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
  it("can expand the filters which are hidden by default", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get(':nth-child(1) > .govuk-checkboxes__item > .govuk-label').should('contain', 'Casework')
  })

  it("can hide the filters", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('#option-select-title-task-type').click()
    cy.get('.govuk-fieldset > :nth-child(1) > .govuk-checkboxes__item > .govuk-label').should('not.be.visible')
  })

  it("can apply a filter which adds task type heading", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[type="checkbox"]').eq(0).check()
    cy.get('#actionFilter').click()
    cy.url().should('include', 'selected-task-type=CWGN')
    cy.get('.moj-filter__selected').should('contain','Task type')
  })

  it("can apply two filters", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[type="checkbox"]').eq(0).check()
    cy.get('[type="checkbox"]').eq(1).check()
    cy.get('#actionFilter').click()
    cy.url().should('include', 'selected-task-type=CWGN')
    cy.url().should('include', 'selected-task-type=ORAL')
  })
  
  it("retains task type filter when changing views", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[type="checkbox"]').eq(0).check()
    cy.get('#actionFilter').click()
    cy.get("#display-rows").select('100')
    cy.url().should('include', 'selected-task-type=CWGN')
  })

  it("shows button to remove individual task type filter", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[type="checkbox"]').eq(0).check()
    cy.get('[type="checkbox"]').eq(1).check()
    cy.get('#actionFilter').click()
    cy.get('.moj-filter__tag').eq(0).should('contain', 'Casework - General')
    cy.get('.moj-filter__tag').eq(1).should('contain', 'Order - Allocate to team')
  })

  it("can clear all filters with clear filter link", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[type="checkbox"]').eq(0).check()
    cy.get('[type="checkbox"]').eq(1).check()
    cy.get('#actionFilter').click()
    cy.get('#clear-filters').click()
    cy.get('.moj-filter__tag').should('not.exist');
    cy.get('[type="checkbox"]').eq(0).should('not.be.checked') 
    cy.get('[type="checkbox"]').eq(1).should('not.be.checked') 
  })

})