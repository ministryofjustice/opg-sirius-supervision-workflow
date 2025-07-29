describe("Footer", () => {
    beforeEach(() => {
        cy.visit('/caseload?team=21');
    });

    it("should show the accessibility link", () => {
        cy.get('[data-cy="accessibilityStatement"]').invoke("removeAttr", "target").click();
        cy.on("url:changed", (newUrl) => {
            expect(newUrl).to.contain("accessibility");
        });
    });
});
