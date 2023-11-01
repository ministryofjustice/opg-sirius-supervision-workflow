describe("Filters", () => {

  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.window().then((win) => {
      win.sessionStorage.clear()
    })
    cy.visit("/deputies?team=27");
});
  it("can expand the filters which are hidden by default", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('#list-of-tasks-to-filter label').should('contain', 'Casework')
  })

  it("can hide the filters", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('#option-select-title-task-type').click()
    cy.get('#list-of-tasks-to-filter label').should('not.be.visible')
  })

  it("can apply a filter which adds assignee heading", () => {
    cy.get('#option-select-title-assignee').click()
    cy.get('[data-filter-name="moj-filter-name-assignee"]').within(() => {
      cy.get('label:contains("Not Assigned")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'unassigned=21')
    cy.get('.moj-filter__selected').should('contain','Case owner')
  })
})
