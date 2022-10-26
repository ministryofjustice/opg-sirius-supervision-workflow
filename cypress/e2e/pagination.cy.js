describe("Pagination", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/workflow/1");
  });

  describe("View 25", () => {
    it("allows me to select view 25 and updates task numbers", () => {
      cy.get("#display-rows").select('25')
      cy.get("#display-rows").should('have.value', '25')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '1')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '25')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '101')
      cy.get("#next-page-pagination-link").click()
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '26')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '50')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '101')
      cy.get("#next-page-pagination-link").click()
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '51')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '75')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '101')
      cy.get("#next-page-pagination-link").click()
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '76')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '100')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '101')
      cy.get("#next-page-pagination-link").click()
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '101')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '101')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '101')
    })

    it("will not show previous arrow on page 1 but will show it on other pages", () => {
      cy.get("#display-rows").select('25')
      cy.get("#display-rows").should('have.value', '25')
      cy.get("#previous-page-pagination-link").should('not.be.visible', 'Previous')
      cy.get("#next-page-pagination-link").click()
      cy.get("#previous-page-pagination-link").should('be.visible', 'Previous')
    })

    it("shows next button apart from one last page", () => {
      cy.get("#next-page-pagination-link").should('be.visible', 'Next')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .govuk-pagination__list > :nth-child(5) > .govuk-link").click()
      cy.get("#next-page-pagination-link").should('not.be.visible', 'Next')
    })

    it("shows first page and ellipses once you are past page three", () => {
      cy.get("#display-rows").select('25')
      cy.get("#display-rows").should('have.value', '25')
      cy.get("#first-ellipses").should('not.exist');
      cy.get("#first-page-pagination-link").should('not.exist');
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .govuk-pagination__next > .govuk-link").click()
      cy.get("#first-ellipses").should('not.exist');
      cy.get("#first-page-pagination-link").should('not.exist');
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .govuk-pagination__next > .govuk-link").click()
      cy.get("#first-ellipses").should('not.exist');
      cy.get("#first-page-pagination-link").should('not.exist');
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .govuk-pagination__next > .govuk-link").click()
      cy.get("#first-ellipses").should('exist');
      cy.get("#first-page-pagination-link").should('exist');
    })

    it("shows first page and ellipses once you are past page three", () => {
      cy.get("#display-rows").select('25')
      cy.get("#display-rows").should('have.value', '25')
      cy.get("#first-ellipses").should('not.exist');
      cy.get("#next-page-pagination-link").click()
      cy.get("#first-ellipses").should('not.exist');
      cy.get("#next-page-pagination-link").click()
      cy.get("#first-ellipses").should('not.exist');
      cy.get("#next-page-pagination-link").click()
      cy.get("#first-ellipses").should('exist');
    })
  });

  describe("View 50", () => {
    it("can select 50 from task view value dropdown and correctly show when one task on a page", () => {
      cy.get("#display-rows").select('50')
      cy.get("#display-rows").should('have.value', '50')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '1')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '50')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '51')
      cy.get("#next-page-pagination-link").click()
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '51')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '51')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '51')
    })
  });

  describe("View 100", () => {
    it("can select 100 from task view value dropdown and shows limited task count", () => {
      cy.get("#display-rows").select('100')
      cy.get("#display-rows").should('have.value', '100')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '1')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '10')
      cy.get(".govuk-\\!-padding-right-2 > #pagination-label > .flex-container > .moj-pagination__results").should('contain', '10')
    })
  });
});