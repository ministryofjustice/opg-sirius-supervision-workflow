describe("Edit my details", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/my-details/edit");
    });

    it("shows my phone number", () => {
        cy.get("#f-phonenumber").should("have.value", "03004560300");
    });

    it("allows me to change my phone number", () => {
        cy.get("#f-phonenumber").clear().type("123456789");
        cy.get("button[type=submit]").click();

        cy.contains(".moj-banner", "You have successfully edited your details.");
    });
});
