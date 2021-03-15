describe("Add team member", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/teams/add-member/65");
    });

    it("allows me to add a user to a team", () => {
        cy.get(".govuk-table").should("not.exist");

        cy.get("#f-search").clear().type("admin");
        cy.get("button[type=submit]").click();

        cy.get(".govuk-table__row").should("have.length", 2);

        const expected = [
            "system admin",
            "system.admin@opgtest.com",
            "Add to team",
        ];

        cy.get(".govuk-table__body > .govuk-table__row")
            .children()
            .each(($el, index) => {
                cy.wrap($el).should("contain", expected[index]);
            });

        cy.contains("button", "Add to team").click();

        cy.contains(
            ".moj-banner",
            "You have successfully added system.admin@opgtest.com to the team."
        );
    });
});
