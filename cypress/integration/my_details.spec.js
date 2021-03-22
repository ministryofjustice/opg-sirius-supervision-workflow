describe("Work flow", () => {
    beforeEach(() => {
        cy.visit("/supervision/workflow", {
            headers: {
                Cookie: "XSRF-TOKEN=abcde; Other=other",
                "OPG-Bypass-Membrane": "1",
                "X-XSRF-TOKEN": "abcde",
            },
        });
    });

    it("shows task", () => {
        const expected = [
            ["Casework - General"],
        ];

        cy.get("#hook-task-type").each(($el, index) => {
            cy.wrap($el).within(() => {
                cy.get(".govuk-checkboxes__label").should(
                    "have.text",
                    expected[index][0]
                );
            });
        });
    });
});
