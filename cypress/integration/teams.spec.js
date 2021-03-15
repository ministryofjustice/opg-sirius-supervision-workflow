describe("Teams", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/teams");
    });

    it("lists all teams", () => {
        cy.get(".govuk-table__row").should("have.length", 2);

        const expected = ["Cool Team", "Supervision — Allocations", "1"];

        cy.get(".govuk-table__body > .govuk-table__row")
            .children()
            .each(($el, index) => {
                cy.wrap($el).should("contain", expected[index]);
            });
    });

    it("allows me to search for a team", () => {
        cy.get("#f-search").clear().type("cool");
        cy.get("button[type=submit]").click();

        cy.get(".govuk-table__body > .govuk-table__row").should("have.length", 1);

        cy.get("#f-search").clear().type("no such team");
        cy.get("button[type=submit]").click();

        cy.get(".govuk-table__body > .govuk-table__row").should("have.length", 0);
    });

    it("allows me to add a new team", () => {
        cy.contains(".govuk-button", "Add new team");
    });
});
