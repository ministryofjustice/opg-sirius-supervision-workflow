describe("Work flow", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/workflow");
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
