describe("Filters", () => {

  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.window().then((win) => {
      win.sessionStorage.clear();
    });
    cy.visit("/deputies?team=27");
});
  it("applies and removes an ECM filter", () => {
    cy.get('[data-filter-name="moj-filter-name-ecm"]').click();
      cy.get('[data-filter-name="moj-filter-name-ecm"]')
        .find('label:contains("PROTeam1 User1")').click();
    cy.get('[data-module=apply-filters]').click();

    cy.url().should('include', 'ecm=96');
    cy.get('.moj-filter__selected').should('contain','ECM');
  });
});
