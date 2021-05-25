describe("Reassign Tasks", () => {

  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.visit("/");
});
  it("can expand the filters which are hidden by default", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get(':nth-child(2) > .govuk-checkboxes__item > .govuk-label').should('contain', 'Casework')
  })

  it("can hide the filters", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('#option-select-title-task-type').click()
    cy.get('app-c-option-select__container-inner').should('not.be.visible') 
  })

  it("can apply a filter", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get(':nth-child(2) > .govuk-checkboxes__item > #input-element-tasktype-filter').click()
    cy.get('#actionFilter').click()
    cy.url().should('include', 'selected-task-type=CWGN')
  })

  it("can apply two filters", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get(':nth-child(2) > .govuk-checkboxes__item > #input-element-tasktype-filter').click()
    cy.get(':nth-child(3) > .govuk-checkboxes__item > #input-element-tasktype-filter').click()
    cy.get('#actionFilter').click()
    cy.url().should('include', 'selected-task-type=CWGN&selected-task-type=ORAL')
  })
  
  it("retains task type filter when changing views", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get(':nth-child(2) > .govuk-checkboxes__item > #input-element-tasktype-filter').click()
    cy.get('#actionFilter').click()
    cy.get("#display-rows").select('100')
    cy.url().should('include', 'selected-task-type=CWGN')
  })

})