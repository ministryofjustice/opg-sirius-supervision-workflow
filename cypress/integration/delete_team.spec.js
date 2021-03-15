describe("Delete a team", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/teams/delete/65");
    });

    it("shows the team details", () => {
        cy.get(".govuk-body").should("contain", "Are you sure you want to delete the team Cool Team?");
    });

    it("allows me to delete the team", () => {
        cy.contains("button", "Delete team").click();
        cy.url().should("include", "/teams");
    });
});
