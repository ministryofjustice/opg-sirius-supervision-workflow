describe("Filters", () => {

  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.window().then((win) => {
      win.sessionStorage.clear()
    })
    cy.visit("/caseload?team=21");
});
  it("can expand the filters which are hidden by default", () => {
    cy.get('#option-select-title-assignee').click()
    cy.get('#list-of-assignees-to-filter label').should('contain', 'Not Assigned')
  })

  it("can hide the filters", () => {
    cy.get('#option-select-title-assignee').click()
    cy.get('#list-of-assignees-to-filter label').should('be.visible')
    cy.get('#option-select-title-assignee').click()
    cy.get('#list-of-assignees-to-filter label').should('not.be.visible')
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

  it("can apply two filters", () => {
    cy.get('#option-select-title-assignee').click()
    cy.get('[data-filter-name="moj-filter-name-assignee"]').within(() => {
      cy.get('label:contains("Not Assigned")').click()
      cy.get('label:contains("LayTeam1 User1")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'unassigned=21')
    cy.url().should('include', 'assignee=766')
  })

  it("shows button to remove individual assignee filter", () => {
    cy.get('#option-select-title-assignee').click()
    cy.get('[data-filter-name="moj-filter-name-assignee"]').within(() => {
      cy.get('label:contains("Not Assigned")').click()
      cy.get('label:contains("LayTeam1 User1")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.get('.moj-filter__tag').eq(0).should('contain', 'Not Assigned')
    cy.get('.moj-filter__tag').eq(1).should('contain', 'LayTeam1 User1')
  })

  it("can clear all filters with clear filter link", () => {
    cy.get('#option-select-title-assignee').click()
    cy.get('[data-filter-name="moj-filter-name-assignee"]').within(() => {
      cy.get('label:contains("Not Assigned")').click()
      cy.get('label:contains("LayTeam1 User1")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.get('[data-module=clear-filters]').click()
    cy.get('.moj-filter__tag').should('not.exist');
    cy.get('[type="checkbox"]').should('not.be.checked')
  })
})
