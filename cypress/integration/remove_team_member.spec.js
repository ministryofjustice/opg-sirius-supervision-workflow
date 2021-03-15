describe("Team", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/teams/65");
    });

    it("allows me to remove a member", () => {
        cy.get("input[type=checkbox]").eq(0).check();
        cy.get("button[type=submit]").click();

        cy.url().should("include", "/teams/remove-member/65");
        cy.get(".govuk-body").should("contain", "Are you sure you want to remove John from the Cool Team team?");

        cy.get("button[type=submit]").click();
        cy.url().should("include", "/teams/65");
    });
});
