describe("My details", () => {
    beforeEach(() => {
        cy.visit("/my-details", {
            headers: {
                Cookie: "XSRF-TOKEN=abcde; Other=other",
                "OPG-Bypass-Membrane": "1",
                "X-XSRF-TOKEN": "abcde",
            },
        });
    });

    it("shows my details", () => {
        const expected = [
            ["Name", "system admin"],
            ["Email", "system.admin@opgtest.com"],
            ["Phone number", "03004560300"],
            ["Organisation", ""],
            ["Team", "Allocations - (Supervision)"],
            ["Roles", "System Admin"],
        ];

        cy.get(".govuk-summary-list__row").each(($el, index) => {
            cy.wrap($el).within(() => {
                cy.get(".govuk-summary-list__key").should(
                    "have.text",
                    expected[index][0]
                );
                cy.get(".govuk-summary-list__value").should(
                    "have.text",
                    expected[index][1]
                );
            });
        });
    });

    it("allows me to edit my phone number", () => {
        cy.contains(".govuk-link", "Change phone number");
    });
});
