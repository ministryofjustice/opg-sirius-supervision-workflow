describe("Pagination", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/workflow/1");
  });

  it("shows correct number of total tasks", () => {
    cy.get("#pagination-label > .flex-container > .moj-pagination__results > :nth-child(1)").should('contain', '1')
    cy.get("#pagination-label > .flex-container > .moj-pagination__results > :nth-child(2)").should('contain', '10')
    cy.get("#pagination-label > .flex-container > .moj-pagination__results > :nth-child(3)").should('contain', '10')
  })

 it("disabled previous button while on top page one", () => {
    cy.get("#pagination-label > .flex-container > .moj-pagination__list > .moj-pagination__item--prev > .moj-pagination__link").should('be.disabled')
  })

  it("can select 25 from task view value dropdown", () => {
    cy.get("#display-rows").select('25')
    cy.get("#display-rows").should('have.value', '25')
  })

  it("can select 50 from task view value dropdown", () => {
    cy.get("#display-rows").select('50')
    cy.get("#display-rows").should('have.value', '50')
  })

  it("can select 100 from task view value dropdown", () => {
    cy.get("#display-rows").select('100')
    cy.get("#display-rows").should('have.value', '100')
  })

});