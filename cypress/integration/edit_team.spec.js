describe("Edit a team", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/teams/edit/65");
    });

    it("shows the team details", () => {
        cy.get("#f-name").should("have.value", "Cool Team");
        cy.get("#f-service-conditional").should("be.checked");
        cy.get("#f-service-conditional-2").should("not.be.checked");
        cy.get("#f-type").should("have.value", "ALLOCATIONS");
        cy.get("#f-phoneNumber").should("have.value", "01818118181");
        cy.get("#f-email").should("have.value", "coolteam@opgtest.com");
    });

    it("allows me to change the team's details", () => {
        cy.get("#f-name").clear().type("Another team");
        cy.get("#f-type").select("ALLOCATIONS");
        cy.get("#f-phoneNumber").clear().type("03573953");
        cy.get("#f-email").clear().type("other.team@opgtest.com");
        cy.get("button[type=submit]").click();

        cy.contains(".moj-banner", "You have successfully edited Another team.");
    });
});
