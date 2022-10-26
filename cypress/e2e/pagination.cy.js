describe("Pagination", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/workflow/1");
  });

  // it("does not show previous button on page one", () => {
  //   cy.get(".govuk-\\!-padding-right-2 > #pagination > .flex-container").should('not.be.visible', 'Previous')
  // })

  it("shows next button on page one", () => {
    cy.get(".govuk-\\!-padding-right-2 > #pagination > .flex-container").should('contain', 'Next')
  })

  it("can select 25 from task view value dropdown", () => {
    cy.get("#display-rows").select('25')
    cy.get("#display-rows").should('have.value', '25')
    cy.get(".govuk-\\!-padding-right-2 > #pagination > .flex-container").should('contain', '1')
    cy.get(".govuk-\\!-padding-right-2 > #pagination > .flex-container").should('contain', '25')
    cy.get(".govuk-\\!-padding-right-2 > #pagination > .flex-container").should('contain', '101')
  })

  it("can select 50 from task view value dropdown", () => {
    cy.get("#display-rows").select('50')
    cy.get("#display-rows").should('have.value', '50')
    cy.get(".govuk-\\!-padding-right-2 > #pagination > .flex-container").should('contain', '1')
    cy.get(".govuk-\\!-padding-right-2 > #pagination > .flex-container").should('contain', '50')
    cy.get(".govuk-\\!-padding-right-2 > #pagination > .flex-container").should('contain', '101')
  })

  it("can select 100 from task view value dropdown", () => {
    cy.get("#display-rows").select('100')
    cy.get("#display-rows").should('have.value', '100')
    cy.get(".govuk-\\!-padding-right-2 > #pagination > .flex-container").should('contain', '1')
    cy.get(".govuk-\\!-padding-right-2 > #pagination > .flex-container").should('contain', '100')
    cy.get(".govuk-\\!-padding-right-2 > #pagination > .flex-container").should('contain', '101')
  })

});