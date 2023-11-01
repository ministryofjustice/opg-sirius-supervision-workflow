describe("Filters", () => {

  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.window().then((win) => {
      win.sessionStorage.clear()
    })
    cy.visit("/deputy-tasks?team=27");
  });

  it("can expand the filters which are hidden by default", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('#list-of-tasks-to-filter label').should('contain', 'PDR follow up')
  })

  it("can hide the filters", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('#option-select-title-task-type').click()
    cy.get('#list-of-tasks-to-filter label').should('not.be.visible')
  })

  it("can apply a filter which adds task type heading", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[data-filter-name="moj-filter-name-tasktype"]').within(() => {
      cy.get('label:contains("PDR follow up")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'task-type=PFU')
    cy.get('.moj-filter__selected').should('contain','Task type')
  })

  it("can apply a filter which adds assignee heading", () => {
    cy.get('#option-select-title-assignee').click()
    cy.get('[data-filter-name="moj-filter-name-assignee"]').within(() => {
      cy.get('label:contains("Not Assigned")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'unassigned=27')
    cy.get('.moj-filter__selected').should('contain','Not Assigned')
  })

  it("can apply all filters", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[data-filter-name="moj-filter-name-tasktype"]').within(() => {
      cy.get('label:contains("PDR follow up")').click()
      cy.get('label:contains("Quarterly catch up call")').click()
    })
    cy.get('#option-select-title-assignee').click()
    cy.get('[data-filter-name="moj-filter-name-assignee"]').within(() => {
      cy.get('label:contains("Not Assigned")').click()
      cy.get('label:contains("PROTeam1 User1")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'task-type=PFU')
    cy.url().should('include', 'task-type=QCUC')
    cy.url().should('include', 'unassigned=27')
    cy.url().should('include', 'assignee=96')
  })

  // it("retains task type filter when changing views", () => {
  //   cy.get('#option-select-title-task-type').click()
  //   cy.get('[data-filter-name="moj-filter-name-tasktype"]').within(() => {
  //     cy.get('label:contains("PDR follow up")').click()
  //   })
  //   cy.get('[data-module=apply-filters]').click()
  //   cy.get("#top-pagination .display-rows").select('100')
  //   cy.url().should('include', 'task-type=PFU')
  // })

  it("shows button to remove individual task type filter", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[data-filter-name="moj-filter-name-tasktype"]').within(() => {
      cy.get('label:contains("PDR follow up")').click()
      cy.get('label:contains("Quarterly catch up call")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.get('.moj-filter__tag').eq(0).should('contain', 'PDR follow up')
    cy.get('.moj-filter__tag').eq(1).should('contain', 'Quarterly catch up call')
  })

  it("shows button to remove individual assignee filter", () => {
    cy.get('#option-select-title-assignee').click()
    cy.get('[data-filter-name="moj-filter-name-assignee"]').within(() => {
      cy.get('label:contains("Not Assigned")').click()
      cy.get('label:contains("PROTeam1 User1")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.get('.moj-filter__tag').eq(0).should('contain', 'Not Assigned')
    cy.get('.moj-filter__tag').eq(1).should('contain', 'PROTeam1 User1')
  })

  it("can clear all filters with clear filter link", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[data-filter-name="moj-filter-name-tasktype"]').within(() => {
      cy.get('label:contains("PDR follow up")').click()
      cy.get('label:contains("Quarterly catch up call")').click()
    })
    cy.get('#option-select-title-assignee').click()
    cy.get('[data-filter-name="moj-filter-name-assignee"]').within(() => {
      cy.get('label:contains("Not Assigned")').click()
      cy.get('label:contains("PROTeam1 User1")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.get('.moj-filter__tag').should('be.visible');
    cy.get('[data-module=clear-filters]').click()
    cy.get('.moj-filter__tag').should('not.exist');
    cy.get('[type="checkbox"]').should('not.be.checked')
  })
})
