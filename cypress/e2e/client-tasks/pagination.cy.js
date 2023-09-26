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
      cy.get(".moj-team-banner__container > .govuk-form-group > .govuk-select").select("Lay Team 2 - (Supervision)")
    })

    it("will not show previous arrow on page 1 but will show it on other pages", () => {
      cy.get('#top-nav > :nth-child(1) > .flex-container > .govuk-pagination__list > :nth-child(2) > .govuk-link').click({force: true})
      cy.get(".previous-page-pagination-link").should('exist')
      cy.get('#top-nav > :nth-child(1) > .flex-container > .govuk-pagination__list > :nth-child(3) > .govuk-link').click({force: true})
      cy.get(".previous-page-pagination-link").should('be.visible', 'Previous')
      cy.get('#top-nav > :nth-child(1) > .flex-container > .govuk-pagination__list > :nth-child(4) > .govuk-link').click({force: true})
      cy.get(".previous-page-pagination-link").should('be.visible', 'Previous')
      cy.get('#top-nav > :nth-child(1) > .flex-container > .govuk-pagination__list > :nth-child(5) > .govuk-link').click({force: true})
    })

    it("shows next button apart from on last page", () => {
      cy.get("#top-nav .next-page-pagination-link").should('be.visible', 'Next')
      cy.get('#top-nav > :nth-child(1) > .flex-container > .govuk-pagination__list > :nth-child(2) > .govuk-link').click({force: true})
      cy.get("#top-nav .next-page-pagination-link").should('be.visible', 'Next')
      cy.get('#top-nav > :nth-child(1) > .flex-container > .govuk-pagination__list > :nth-child(3) > .govuk-link').click({force: true})
      cy.get("#top-nav .next-page-pagination-link").should('be.visible', 'Next')
      cy.get('#top-nav > :nth-child(1) > .flex-container > .govuk-pagination__list > :nth-child(4) > .govuk-link').click({force: true})
      cy.get("#top-nav .next-page-pagination-link").should('be.visible', 'Next')
      cy.get('#top-nav > :nth-child(1) > .flex-container > .govuk-pagination__list > :nth-child(5) > .govuk-link').click({force: true})
      cy.get("#top-nav .next-page-pagination-link").should('not.exist')
    })

    it("shows first and second ellipsis when expected", () => {
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
      cy.get(".moj-team-banner__container > .govuk-form-group > .govuk-select").select("Lay Team 2 - (Supervision)")
      cy.get(".moj-pagination__results").should('contain', '1')
      cy.get(".moj-pagination__results").should('contain', '25')
      cy.get(".moj-pagination__results").should('contain', '101')
    })
  });

  describe("View 50", () => {
    beforeEach(() => {
      cy.get("#top-nav .display-rows").select('50')
          .invoke('val').should('contain', 'per-page=50')
      cy.get(".moj-team-banner__container > .govuk-form-group > .govuk-select").select("Lay Team 2 - (Supervision)")
    })

    it("can select 50 from task view value dropdown and correctly show when one task on a page", () => {
      cy.get(".moj-pagination__results").should('contain', '1')
      cy.get(".moj-pagination__results").should('contain', '50')
      cy.get(".moj-pagination__results").should('contain', '51')
    })

  });

  describe("View 100", () => {
    beforeEach(() => {
      cy.get("#top-nav .display-rows").select('100')
          .invoke('val').should('contain', 'per-page=100')
      cy.get(".moj-team-banner__container > .govuk-form-group > .govuk-select").select("Lay Team 2 - (Supervision)")
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