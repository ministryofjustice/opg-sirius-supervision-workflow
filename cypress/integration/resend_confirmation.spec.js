describe("Resend confirmation", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/edit-user/123");
    });

    it("allows me to send an activation email", () => {
        cy.contains(".moj-banner", "User has not activated their account yet.");
        cy.contains("button", "Resend activation email").click();

        cy.url().should("include", "/resend-confirmation");
        cy.contains(
            "A new activation email has been sent to system.admin@opgtest.com"
        );

        cy.contains("a", "Continue").click();

        cy.url().should("include", "/edit-user/123");
    });
});
