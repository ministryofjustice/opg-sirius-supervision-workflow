describe("Unlock user", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/unlock-user/123");
    });

    it("allows me to unlock a user", () => {
        cy.contains("button", "Unlock account").click();
        cy.url().should("include", "/edit-user/123");
    });
});
