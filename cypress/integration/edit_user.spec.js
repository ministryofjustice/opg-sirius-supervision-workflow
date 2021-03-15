describe("Edit user", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/edit-user/123");
    });

    it("allows me to edit a user", () => {
        cy.get("#f-firstname").type("2");
        cy.get("button[type=submit]").click();

        cy.contains(".moj-banner", "You have successfully edited a user.");
    });
});
