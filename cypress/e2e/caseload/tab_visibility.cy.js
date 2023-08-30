describe("Caseload visibility", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    })

    it("Caseload tab is visible and clickable for Lay teams", () => {
        cy.visit("/client-tasks?team=21");

        cy.get(".moj-sub-navigation__item:nth-child(1) a:contains('Client tasks')").as("tab1")
        cy.get(".moj-sub-navigation__item:nth-child(2) a:contains('Caseload')").as("tab2")

        cy.get("@tab1").should("have.attr", "aria-current", "page")
        cy.get("@tab1").should("not.have.attr", "href")

        cy.get("@tab2").should("not.have.attr", "aria-current")
        cy.get("@tab2").should("have.attr", "href", "caseload?team=21")
        cy.get("@tab2").click()

        cy.url().should('contain', '/caseload?team=21')

        cy.get("@tab1").should("not.have.attr", "aria-current")
        cy.get("@tab1").should("have.attr", "href", "client-tasks?team=21")

        cy.get("@tab2").should("have.attr", "aria-current", "page")
        cy.get("@tab2").should("not.have.attr", "href")
    })

    it("Caseload tab is not visible for non-Lay teams", () => {
        cy.visit("/client-tasks?team=13");

        cy.get(".moj-sub-navigation__item:nth-child(1) a:contains('Client tasks')").as("tab1")
        cy.get(".moj-sub-navigation__item:nth-child(2) a:contains('Caseload')").should("not.exist")
    })
})
