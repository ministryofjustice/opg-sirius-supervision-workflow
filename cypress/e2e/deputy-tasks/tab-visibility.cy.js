describe("Deputy tasks visibility", () => {
    before(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    })

    let teams = [24, 27] // PA and Pro
    teams.forEach((team) => {
        it("Deputy tasks tab is visible and clickable for Pro/PA teams", () => {
            cy.visit("/client-tasks?team="+team);

            cy.get(".moj-sub-navigation__item:nth-child(1) a:contains('Client tasks')").as("tab1")
            cy.get(".moj-sub-navigation__item:nth-child(2) a:contains('Deputy tasks')").as("tab2")

            cy.get("@tab1").should("have.attr", "aria-current", "page")
            cy.get("@tab1").should("not.have.attr", "href")

            cy.get("@tab2").should("not.have.attr", "aria-current")
            cy.get("@tab2").should("have.attr", "href", "deputy-tasks?team="+team)
            cy.get("@tab2").click()

            cy.url().should('contain', '/deputy-tasks?team='+team)

            cy.get("@tab1").should("not.have.attr", "aria-current")
            cy.get("@tab1").should("have.attr", "href", "client-tasks?team="+team)

            cy.get("@tab2").should("have.attr", "aria-current", "page")
            cy.get("@tab2").should("not.have.attr", "href")
        })
    });

    it("Deputy tasks tab is not visible for Lay teams", () => {
        cy.visit("/client-tasks?team=21");

        cy.get(".moj-sub-navigation__item:nth-child(1) a:contains('Client tasks')").as("tab1")
        cy.get(".moj-sub-navigation__item a:contains('Deputy tasks')").should("not.exist")
    })
})
