describe("Pagination", () => {
    before(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    let assertPaginationHasLoaded = function () {
        cy.get("#top-pagination").should("exist");
        cy.get("#bottom-pagination").should("exist");
        cy.get(".moj-pagination__results").should("contain.text", "Showing 1 to 2 of 2 deputies")
        cy.get(".govuk-pagination__item:nth-child(1)").should("have.length", 2)
        cy.get(".govuk-pagination__item:nth-child(2)").should("not.exist")
    }

    it("is visible on the Deputies list page", () => {
        cy.visit("/deputies?team=27");
        assertPaginationHasLoaded();
    })
});
