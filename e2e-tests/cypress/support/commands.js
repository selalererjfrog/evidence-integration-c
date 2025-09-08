// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************

// Custom command to wait for service to be ready
Cypress.Commands.add('waitForService', (serviceUrl, endpoint, maxRetries = 30) => {
  let attempts = 0
  
  const checkService = () => {
    cy.request({
      method: 'GET',
      url: `${serviceUrl}${endpoint}`,
      failOnStatusCode: false,
      timeout: 5000
    }).then((response) => {
      if (response.status === 200) {
        cy.log(`✅ Service at ${serviceUrl}${endpoint} is ready`)
      } else {
        attempts++
        if (attempts < maxRetries) {
          cy.wait(2000)
          checkService()
        } else {
          throw new Error(`Service at ${serviceUrl}${endpoint} failed to start after ${maxRetries} attempts`)
        }
      }
    })
  }
  
  checkService()
})
