describe('Quote Service E2E Tests', () => {
  const quoteServiceUrl = Cypress.env('quoteServiceUrl')

  beforeEach(() => {
    // Wait for quote service to be ready
    cy.waitForService(quoteServiceUrl, '/actuator/health')
  })

  it('should have healthy quote service', () => {
    cy.request('GET', `${quoteServiceUrl}/actuator/health`)
      .then((response) => {
        expect(response.status).to.eq(200)
        expect(response.body.status).to.eq('UP')
      })
  })

  it('should return health status from API endpoint', () => {
    cy.request('GET', `${quoteServiceUrl}/api/quotes/health`)
      .then((response) => {
        expect(response.status).to.eq(200)
      })
  })

  it('should return today\'s quote', () => {
    cy.request('GET', `${quoteServiceUrl}/api/quotes/today`)
      .then((response) => {
        expect(response.status).to.eq(200)
        expect(response.body).to.have.property('text')
        expect(response.body).to.have.property('author')
        expect(response.body.text).to.be.a('string')
        expect(response.body.author).to.be.a('string')
      })
  })
})
