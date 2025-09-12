describe("Pagination", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    let assertPaginationHasLoaded = function () {
        cy.get("#top-pagination").should("exist");
        cy.get("#bottom-pagination").should("exist");
        cy.get(".moj-pagination__results").should("contain.text", "Showing 1 to 1 of 1 clients")
        cy.get(".govuk-pagination__item:nth-child(1)").should("have.length", 2)
        cy.get(".govuk-pagination__item:nth-child(2)").should("not.exist")
    }

    it("is visible on the Caseload list page", () => {
        cy.visit("/caseload?team=21");
        assertPaginationHasLoaded();
    })

    it("is visible on the New Deputy Orders list page", () => {
        cy.visit("/caseload?team=28");
        assertPaginationHasLoaded();
    })

    it("is visible on the Closed Clients list page", () => {
        cy.visit("/caseload?team=40");
        assertPaginationHasLoaded();
    })
});
