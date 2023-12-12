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
    cy.get('#option-select-title-status').click()
    cy.get('#list-of-statuses-to-filter label').should('contain', 'Active')
    cy.get('#list-of-statuses-to-filter label').should('contain', 'Closed')
  })

  it("can hide the filters", () => {
    cy.get('#option-select-title-assignee').click()
    cy.get('#list-of-assignees-to-filter label').should('be.visible')
    cy.get('#option-select-title-assignee').click()
    cy.get('#list-of-assignees-to-filter label').should('not.be.visible')

    cy.get('#option-select-title-status').click()
    cy.get('#list-of-statuses-to-filter label').should('be.visible')
    cy.get('#option-select-title-status').click()
    cy.get('#list-of-statuses-to-filter label').should('not.be.visible')
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

  it("can apply a filter which adds status heading", () => {
    cy.get('#option-select-title-status').click()
    cy.get('[data-filter-name="moj-filter-name-status"]').within(() => {
      cy.get('label:contains("Closed")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'status=closed')
    cy.get('.moj-filter__selected').should('contain','Closed')
  })

  it("can apply all filters", () => {
    cy.get('#option-select-title-assignee').click()
    cy.get('[data-filter-name="moj-filter-name-assignee"]').within(() => {
      cy.get('label:contains("Not Assigned")').click()
      cy.get('label:contains("LayTeam1 User1")').click()
    })
    cy.get('#option-select-title-status').click()
    cy.get('[data-filter-name="moj-filter-name-status"]').within(() => {
      cy.get('label:contains("Active")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'unassigned=21')
    cy.url().should('include', 'assignee=766')
    cy.url().should('include', 'status=active')
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

  it("shows button to remove individual status filter", () => {
    cy.get('#option-select-title-status').click()
    cy.get('[data-filter-name="moj-filter-name-status"]').within(() => {
      cy.get('label:contains("Active")').click()
      cy.get('label:contains("Closed")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.get('.moj-filter__tag').eq(0).should('contain', 'Active')
    cy.get('.moj-filter__tag').eq(1).should('contain', 'Closed')
  })

  it("can clear all filters with clear filter link", () => {
    cy.get('#option-select-title-assignee').click()
    cy.get('[data-filter-name="moj-filter-name-assignee"]').within(() => {
      cy.get('label:contains("Not Assigned")').click()
      cy.get('label:contains("LayTeam1 User1")').click()
    })
    cy.get('#option-select-title-status').click()
    cy.get('[data-filter-name="moj-filter-name-status"]').within(() => {
      cy.get('label:contains("Active")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.get('[data-module=clear-filters]').click()
    cy.get('.moj-filter__tag').should('not.exist');
    cy.get('[type="checkbox"]').should('not.be.checked')
  })

  it("can filter by Deputy Type on the HW Caseload page", () => {
    cy.get('#option-select-title-deputy-type').should("not.exist")
    cy.visit("/caseload?team=29")
    cy.get('#option-select-title-deputy-type').click()
    cy.get('#list-of-deputy-types-to-filter label').should('contain', 'Lay')
    cy.get('#list-of-deputy-types-to-filter label').should('contain', 'Professional')
    cy.get('#list-of-deputy-types-to-filter label').should('contain', 'Public Authority')
    cy.get('[data-filter-name="moj-filter-name-deputy-type"]').within(() => {
      cy.get('label:contains("Lay")').click()
      cy.get('label:contains("Public Authority")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'deputy-type=LAY').and('include', 'deputy-type=PA')
    cy.get('.moj-filter__selected').should('contain','Deputy type')
    cy.get('.moj-filter__tag').should('contain', 'Lay')
    cy.get('.moj-filter__tag').should('contain', 'Public Authority')
    cy.get('[data-module=clear-filters]').click()
    cy.get('.moj-filter__tag').should('not.exist')
    cy.get('[type="checkbox"]').should('not.be.checked')
  })

  it("can filter by Case Type on the HW Caseload page", () => {
    cy.get('#option-select-title-case-type').should("not.exist")
    cy.visit("/caseload?team=29")
    cy.get('#option-select-title-case-type').click()
    cy.get('#list-of-case-types-to-filter label').should('contain', 'Hybrid')
    cy.get('#list-of-case-types-to-filter label').should('contain', 'Dual')
    cy.get('#list-of-case-types-to-filter label').should('contain', 'Health and welfare')
    cy.get('#list-of-case-types-to-filter label').should('contain', 'Property and financial affairs')
    cy.get('[data-filter-name="moj-filter-name-case-type"]').within(() => {
      cy.get('label:contains("Hybrid")').click()
      cy.get('label:contains("Dual")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'case-type=HYBRID').and('include', 'case-type=DUAL')
    cy.get('.moj-filter__selected').should('contain','Case type')
    cy.get('.moj-filter__tag').should('contain', 'Hybrid')
    cy.get('.moj-filter__tag').should('contain', 'Dual')
    cy.get('[data-module=clear-filters]').click()
    cy.get('.moj-filter__tag').should('not.exist')
    cy.get('[type="checkbox"]').should('not.be.checked')
  })

  it("can filter by Supervision Level on the Lay Caseload page", () => {
    cy.get('#option-select-title-supervision-level').click()
    cy.get('#list-of-supervision-levels-to-filter label').should('contain', 'General')
    cy.get('#list-of-supervision-levels-to-filter label').should('contain', 'Minimal')
    cy.get('[data-filter-name="moj-filter-name-supervision-level"]').within(() => {
      cy.get('label:contains("General")').click()
      cy.get('label:contains("Minimal")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'supervision-level=GENERAL').and('include', 'supervision-level=MINIMAL')
    cy.get('.moj-filter__selected').should('contain','Supervision level')
    cy.get('.moj-filter__tag').should('contain', 'General')
    cy.get('.moj-filter__tag').should('contain', 'Minimal')
    cy.get('[data-module=clear-filters]').click()
    cy.get('.moj-filter__tag').should('not.exist')
    cy.get('[type="checkbox"]').should('not.be.checked')
  })
})
