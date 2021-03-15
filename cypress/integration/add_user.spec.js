describe("Add user", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/users");
    });

    it("allows me to add a user", () => {
        cy.contains("a", "Add new user").click();

        cy.get("#f-email").clear().type("123456789");
        cy.get("#f-firstname").clear().type("123456789");
        cy.get("#f-surname").clear().type("123456789");

        cy.get("button[type=submit]").click();

        cy.contains(".moj-banner", "You have successfully added a new user.");
    });
});
