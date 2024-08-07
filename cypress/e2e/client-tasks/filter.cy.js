describe("Filters", () => {

  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.window().then((win) => {
      win.sessionStorage.clear()
    })
    cy.visit("/client-tasks");
  });

  it("includes the task type count", () => {
      cy.contains('ECM Tasks').should('contain', '(21)');
      cy.contains('Casework - Call back request').should('contain', '(3)');
  });

  it("includes the assignee count", () => {
        cy.contains('Allocations User1').should('contain', '(12)');
        cy.contains('Allocations User2').should('contain', '(8)');
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

  it("can apply a filter which adds task type heading", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[data-filter-name="moj-filter-name-tasktype"]').within(() => {
      cy.get('label:contains("Casework - General")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'task-type=CWGN')
    cy.get('.moj-filter__selected').should('contain','Task type')
  })

  it("can apply two filters", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[data-filter-name="moj-filter-name-tasktype"]').within(() => {
      cy.get('label:contains("Casework - General")').click()
      cy.get('label:contains("Order - Allocate to team")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'task-type=CWGN')
    cy.url().should('include', 'task-type=ORAL')
  })

  it("can filter by due date", () => {
    cy.get('#option-select-title-due-date').click()
    cy.get('#due-date-from').type("2022-12-17")
    cy.get('#due-date-to').type("2022-12-18")
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'due-date-from=2022-12-17')
    cy.url().should('include', 'due-date-to=2022-12-18')
    cy.get('.moj-filter__tag').contains('From: 2022-12-17')
    cy.get('.moj-filter__tag').contains('To: 2022-12-18')
    cy.get('#due-date-from').should('have.value', "2022-12-17")
    cy.get('#due-date-to').should('have.value', "2022-12-18")
  })

  it("retains task type filter when changing views", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[data-filter-name="moj-filter-name-tasktype"]').within(() => {
      cy.get('label:contains("Casework - General")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.get("#top-pagination .display-rows").select('100')
    cy.url().should('include', 'task-type=CWGN')
  })

  it("shows button to remove individual task type filter", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[data-filter-name="moj-filter-name-tasktype"]').within(() => {
      cy.get('label:contains("Casework - General")').click()
      cy.get('label:contains("Order - Allocate to team")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.get('.moj-filter__tag').eq(0).should('contain', 'Casework - General')
    cy.get('.moj-filter__tag').eq(1).should('contain', 'Order - Allocate to team')
  })

  it("can clear all filters with clear filter link", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[data-filter-name="moj-filter-name-tasktype"]').within(() => {
      cy.get('label:contains("Casework - General")').click()
      cy.get('label:contains("Order - Allocate to team")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.get('[data-module=clear-filters]').click()
    cy.get('.moj-filter__tag').should('not.exist');
    cy.get('[type="checkbox"]').should('not.be.checked')
  })

  it("shows combined team member filters when viewing combined Lay teams", () => {
    cy.get(".moj-team-banner__container > .govuk-form-group > .govuk-select").select("Lay Deputy Team")
    cy.get(".moj-team-banner__container > h1").should('contain', "Lay Deputy Team")
    cy.get('.govuk-fieldset .govuk-checkboxes__item > .govuk-label').should('contain', "LayTeam1 User1")
    cy.get('.govuk-fieldset .govuk-checkboxes__item > .govuk-label').should('contain', "LayTeam2 User1")
  })

  it("applies the ECM Tasks filter to display only ECM tasks", () => {
    cy.get('#option-select-title-task-type').click()
    cy.get('[data-filter-name="moj-filter-name-tasktype"]').within(() => {
      cy.get('[type="checkbox"]').get('label:contains("ECM Tasks")').click()
    })
    cy.get('[data-module=apply-filters]').click()
    cy.url().should('include', 'task-type=ECM_TASKS')
    cy.get('.moj-filter__tag').eq(0).should('contain', 'ECM Tasks')
  })
})
