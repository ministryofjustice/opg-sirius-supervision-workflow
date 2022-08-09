package support

it('my test', () => {
    cy.once('uncaught:exception', () => false);

    // action that causes exception
    cy.get('body').click();
});
