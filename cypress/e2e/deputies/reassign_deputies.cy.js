describe("Reassign deputies", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/deputies?team=27");
    });

    it("allows deputies to be reassigned", () => {
        cy.intercept('api/v1/teams/27', {
            body: {
                "members": [
                    {
                        "id": 76,
                        "displayName": "ProTeam1 User1",
                        "suspended": false,
                    },
                    {
                        "id": 112,
                        "displayName": "ProTeam1 User2",
                        "suspended": false,
                    },
                ]
            }
        });
        cy.setCookie("success-route", "/reassign-deputies/1");
        cy.get('#select-deputy-13').click();
        cy.get('#manage-deputy').should('be.visible').click();
        cy.get('#assignTeam').select('Pro Team 1 - (Supervision)');
        cy.get('#assignCM').select('ProTeam1 User1');
        cy.intercept('PATCH', 'api/v1/users/*', {statusCode: 204})
        cy.get('#edit-save').click()
        cy.get("#success-banner").should('be.visible')
        cy.get("#success-banner").contains('You have reassigned ')
    })

    it("can cancel out of reassigning", () => {
        cy.get('#manage-deputy').should('not.be.visible');
        cy.get('.moj-manage-list__edit-panel').should('not.be.visible');
        cy.get('#select-deputy-13').click();
        cy.get('#manage-deputy').should('be.visible');
        cy.get('.moj-manage-list__edit-panel').should('not.be.visible');
        cy.get('#manage-deputy').should('be.visible').click();
        cy.get('.moj-manage-list__edit-panel').should('be.visible');
        cy.get('#edit-cancel').click();
        cy.get('#manage-deputy').should('be.visible');
        cy.get('.moj-manage-list__edit-panel').should('not.be.visible');
        cy.get('#select-deputy-13').click();
        cy.get('#manage-deputy').should('not.be.visible');
    })
});
