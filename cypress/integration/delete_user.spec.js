describe.skip("Delete user", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/delete-user/123");
    });

    it("allows me to delete a user", () => {
        cy.contains("button", "Delete account").click();
        cy.url().should("include", "/users");
    });
});
