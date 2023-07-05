describe("Pagination", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.visit("/client-tasks");
  });

  describe("First page, ellipses and final page", () => {
    beforeEach(() => {
      cy.get("#top-nav .display-rows").select('25')
          .invoke('val').should('contain', 'per-page=25')
    })

    it("will not show previous arrow on page 1 but will show it on other pages", () => {
      cy.get(".previous-page-pagination-link").should('not.exist')
      cy.get("#top-nav .next-page-pagination-link").click()
      cy.get(".previous-page-pagination-link").should('be.visible', 'Previous')
    })

    it("shows next button apart from on last page", () => {
      cy.get("#top-nav .next-page-pagination-link").should('be.visible', 'Next')
      cy.get("#top-nav .govuk-pagination__item:nth-last-child(1) .govuk-pagination__link").click()
      cy.get("#top-nav .next-page-pagination-link").should('not.exist')
    })

    it("shows first page and ellipses once you are past page 4", () => {
      let firstEllipsis = ".govuk-pagination__item--ellipses:nth-child(2)",
          secondEllipsis = ".govuk-pagination__item--ellipses:nth-last-child(2)"

      cy.get(firstEllipsis).should("not.exist")
      cy.get(secondEllipsis).should("exist")

      cy.get("#top-nav .govuk-pagination__link:contains(2)").click()

      cy.get(firstEllipsis).should("not.exist")
      cy.get(secondEllipsis).should("not.exist")

      cy.get("#top-nav .govuk-pagination__link:contains(3)").click()

      cy.get(firstEllipsis).should("not.exist")
      cy.get(secondEllipsis).should("not.exist")

      cy.get("#top-nav .govuk-pagination__link:contains(4)").click()

      cy.get(firstEllipsis).should("not.exist")
      cy.get(secondEllipsis).should("not.exist")

      cy.get("#top-nav .govuk-pagination__link:contains(5)").click()

      cy.get(firstEllipsis).should("exist")
      cy.get(secondEllipsis).should("not.exist")
    })
  });

  describe("View 25", () => {
    beforeEach(() => {
      cy.get("#top-nav .display-rows").select('25')
          .invoke('val').should('contain', 'per-page=25')
    })

    it("allows me to select view 25 and updates task numbers", () => {
      cy.get(".moj-pagination__results").should('contain', '1')
      cy.get(".moj-pagination__results").should('contain', '25')
      cy.get(".moj-pagination__results").should('contain', '101')
      cy.get("#top-nav .next-page-pagination-link").click()
      cy.get(".moj-pagination__results").should('contain', '26')
      cy.get(".moj-pagination__results").should('contain', '50')
      cy.get(".moj-pagination__results").should('contain', '101')
      cy.get("#top-nav .next-page-pagination-link").click()
      cy.get(".moj-pagination__results").should('contain', '51')
      cy.get(".moj-pagination__results").should('contain', '75')
      cy.get(".moj-pagination__results").should('contain', '101')
      cy.get("#top-nav .next-page-pagination-link").click()
      cy.get(".moj-pagination__results").should('contain', '76')
      cy.get(".moj-pagination__results").should('contain', '100')
      cy.get(".moj-pagination__results").should('contain', '101')
      cy.get("#top-nav .next-page-pagination-link").click()
      cy.get(".moj-pagination__results").should('contain', '101')
      cy.get(".moj-pagination__results").should('contain', '101')
      cy.get(".moj-pagination__results").should('contain', '101')
    })
  });

  describe("View 50", () => {
    beforeEach(() => {
      cy.get("#top-nav .display-rows").select('50')
          .invoke('val').should('contain', 'per-page=50')
    })

    it("can select 50 from task view value dropdown and correctly show when one task on a page", () => {
      cy.get(".moj-pagination__results").should('contain', '1')
      cy.get(".moj-pagination__results").should('contain', '50')
      cy.get(".moj-pagination__results").should('contain', '51')
      cy.get("#top-nav .next-page-pagination-link").click()
      cy.get(".moj-pagination__results").should('contain', '51')
      cy.get(".moj-pagination__results").should('contain', '51')
      cy.get(".moj-pagination__results").should('contain', '51')
    })

  });

  describe("View 100", () => {
    beforeEach(() => {
      cy.get("#top-nav .display-rows").select('100')
          .invoke('val').should('contain', 'per-page=100')
    })

    it("can select 100 from task view value dropdown and shows limited task count", () => {
      cy.get(".moj-pagination__results").should('contain', '1')
      cy.get(".moj-pagination__results").should('contain', '10')
      cy.get(".moj-pagination__results").should('contain', '10')
    })

    it("will not show previous link, next link, first page, first ellipses or final ellipses when only one page", () => {
      cy.get(".previous-page-pagination-link").should('not.exist');
      cy.get(".next-page-pagination-link").should('not.exist');
      cy.get(".govuk-pagination__item--ellipses").should('not.exist');
      cy.get(".govuk-pagination__item").should('have.length', 2)
      cy.get(".govuk-pagination__item:contains(1)").should('have.length', 2)
    })
  });
});