describe("Team", () => {
    beforeEach(() => {
        cy.visit("/teams/65", {
            headers: {
                Cookie: "XSRF-TOKEN=abcde; Other=other",
                "OPG-Bypass-Membrane": "1",
                "X-XSRF-TOKEN": "abcde",
            },
        });
    });

    it("shows team members", () => {
        cy.get(".govuk-table__row").should("have.length", 2);

        const expected = ["Select", "John", "john@opgtest.com"];

        cy.get(".govuk-table__body > .govuk-table__row")
            .children()
            .each(($el, index) => {
                cy.wrap($el).should("contain", expected[index]);
            });
    });

    it("allows me to edit the team", () => {
        cy.contains(".govuk-button", "Edit team");
    });

    it("allows me to add a team member", () => {
        cy.contains(".govuk-button", "Add user to team");
    });

    it("allows me to remove team members", () => {
        cy.contains(".govuk-button", "Remove selected from team");

        cy.get(".govuk-table__body > .govuk-table__row input[type=checkbox]").should("have.length", 1);
    });
});
