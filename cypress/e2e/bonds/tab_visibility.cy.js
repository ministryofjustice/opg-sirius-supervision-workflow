describe("Bonds visibility", () => {
    before(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    })

    it("Bonds tab is visible and clickable for allocation team", () => {
        cy.visit("/client-tasks?team=13");

        cy.get(".moj-sub-navigation__item:nth-child(1) a:contains('Client tasks')").as("tab1")
        cy.get(".moj-sub-navigation__item:nth-child(2) a:contains('Bonds')").as("tab2")

        cy.get("@tab1").should("have.attr", "aria-current", "page")
        cy.get("@tab1").should("not.have.attr", "href")

        cy.get("@tab2").should("not.have.attr", "aria-current")
        cy.get("@tab2").should("have.attr", "href", "bonds?team=13")
        cy.get("@tab2").click()

        cy.url().should('contain', '/bonds?team=13')

        cy.get("@tab1").should("not.have.attr", "aria-current")
        cy.get("@tab1").should("have.attr", "href", "client-tasks?team=13&preselect")

        cy.get("@tab2").should("have.attr", "aria-current")
        cy.get("@tab2").should("not.have.have.attr", "href")
    });

    it("Bonds tab is not visible for non-allocation teams", () => {
        cy.visit("/client-tasks?team=21");

        cy.get(".moj-sub-navigation__item:nth-child(1) a:contains('Client tasks')").as("tab1")
        cy.get(".moj-sub-navigation__item a:contains('Bonds')").should("not.exist")
    })
})
