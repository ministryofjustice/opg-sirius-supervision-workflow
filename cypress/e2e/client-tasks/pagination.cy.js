describe("Pagination", () => {
  before(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
  });

  it("is visible on the Client tasks list page", () => {
    cy.visit("/client-tasks");
    cy.get("#top-pagination").should("exist");
    cy.get("#bottom-pagination").should("exist");
    cy.get(".moj-pagination__results").should("contain.text", "Showing 1 to 13 of 13 tasks")
    cy.get(".govuk-pagination__item:nth-child(1)").should("have.length", 2)
    cy.get(".govuk-pagination__item:nth-child(2)").should("not.exist")
  })
});
