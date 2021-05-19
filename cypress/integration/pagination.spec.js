describe("Pagination", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/");
  });

  it("shows correct number of total tasks", () => {
    cy.get("#pagination-label > .flex-container > .moj-pagination__results > :nth-child(1)").should('contain', '1')
    cy.get("#pagination-label > .flex-container > .moj-pagination__results > :nth-child(2)").should('contain', '25')
    cy.get("#pagination-label > .flex-container > .moj-pagination__results > :nth-child(3)").should('contain', '100')
  })

  it("changes url when page numbers is clicked", () => {
    cy.get('#pagination-label > .flex-container > .moj-pagination__list > :nth-child(5) > .moj-pagination__link').first().click()
    cy.url().should('include', 'page=2&currentTaskDisplay=25')
  })

  it("changes url when next button is clicked", () => {
    cy.get("#pagination-label > .flex-container > .moj-pagination__list > .moj-pagination__item--next > .moj-pagination__link").first().click()
    cy.url().should('include', 'page=2&currentTaskDisplay=25') 
  })

  it("reloads final page if you are on the final page and press next", () => {
    // get to page 4
    cy.get('#pagination-label > .flex-container > .moj-pagination__list > :nth-child(8) > .moj-pagination__link').first().click()
    cy.get("#pagination-label > .flex-container > .moj-pagination__list > .moj-pagination__item--next > .moj-pagination__link").first().click()
    cy.url().should('include', 'page=4&currentTaskDisplay=25')
  })

  it("changes url when previous button is clicked", () => {
    //cannot currently test previous button as only return one page of tasks (element disabled for page)
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