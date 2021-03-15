describe("Delete user", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/delete-user/123");
    });

    it("allows me to delete a user", () => {
        cy.get(".govuk-body").should("contain", "Are you sure you want to delete system admin?");
        cy.get("button[type=submit]").contains("Delete user").click();

        cy.get('a[href*="/users"]').contains('Continue').click()
        cy.url().should("include", "/users");
    });
});
