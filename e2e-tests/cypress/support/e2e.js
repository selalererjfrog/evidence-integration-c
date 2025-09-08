// ***********************************************************
// This example support/e2e.js is processed and
// loaded automatically before your test files.
//
// This is a great place to put global configuration and
// behavior that modifies Cypress.
//
// You can change the location of this file or turn off
// automatically serving support files with the
// 'supportFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/configuration
// ***********************************************************

// Import commands.js using ES2015 syntax:
import './commands'

// Alternatively you can use CommonJS syntax:
// require('./commands')

// Custom command to wait for service health
Cypress.Commands.add('waitForServiceHealth', (serviceUrl, endpoint, expectedStatus = 'UP') => {
  cy.request({
    method: 'GET',
    url: `${serviceUrl}${endpoint}`,
    failOnStatusCode: false,
    timeout: 30000
  }).then((response) => {
    expect(response.status).to.eq(200)
    if (expectedStatus) {
      expect(response.body.status).to.eq(expectedStatus)
    }
  })
})
