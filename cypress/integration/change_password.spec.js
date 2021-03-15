describe("Change password", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/change-password");
    });

    it("allows me to change my phone number", () => {
        cy.get("#f-currentpassword").clear().type("123456789");
        cy.get("#f-password1").clear().type("123456789");
        cy.get("#f-password2").clear().type("123456789");

        cy.get("button[type=submit]").click();

        cy.contains(".moj-banner", "You have successfully changed your password.");
    });
});
