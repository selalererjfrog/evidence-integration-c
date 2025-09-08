describe('Quotopia End-to-End Tests', () => {
  const quoteServiceUrl = 'http://localhost:8001'
  const uiServiceUrl = 'http://localhost:8081'
  
  before(() => {
    // Check if services are running, if not start them
    cy.task('startServicesIfNeeded', {
      quoteServiceUrl,
      uiServiceUrl
    })
  })

  beforeEach(() => {
    // Wait for both services to be ready
    cy.waitForService(quoteServiceUrl, '/actuator/health')
    cy.waitForService(uiServiceUrl, '/')
  })

  it('should display the Quotopia UI with proper structure', () => {
    cy.visit(uiServiceUrl)
    
    // Check main container structure
    cy.get('[data-testid="quotopia-container"]').should('be.visible')
    cy.get('[data-testid="quotopia-header"]').should('be.visible')
    cy.get('[data-testid="quotopia-main"]').should('be.visible')
    cy.get('[data-testid="quotopia-footer"]').should('be.visible')
    
    // Check banner image
    cy.get('[data-testid="quotopia-banner"]')
      .should('be.visible')
      .and('have.attr', 'src', 'quotopia.png')
      .and('have.attr', 'alt', 'Quotopia Banner')
    
    // Check page title
    cy.title().should('eq', 'Quotopia - Inspire yourself THEN inspire the world')
  })

  it('should display quote loading state initially', () => {
    cy.visit(uiServiceUrl)
    
    // Check that the UI shows loading state
    cy.get('[data-testid="quote-text"]')
      .should('contain', 'Loading today\'s inspiration...')
    
    cy.get('[data-testid="quote-author"]')
      .should('contain', '—')
    
    // Verify quote structure is present
    cy.get('[data-testid="quote-card"]').should('be.visible')
    cy.get('[data-testid="quote-container"]').should('be.visible')
  })

  it('should display quote date element', () => {
    cy.visit(uiServiceUrl)
    
    // Check if date element exists (it might be empty initially)
    cy.get('[data-testid="quote-date"]').should('exist')
  })

  it('should have functional date selector', () => {
    cy.visit(uiServiceUrl)
    
    // Check date selector is present and functional
    cy.get('[data-testid="date-selector"]')
      .should('be.visible')
      .and('have.attr', 'type', 'date')
      .and('have.attr', 'max', '2025-12-31')
    
    // Check label is present
    cy.get('[data-testid="date-label"]')
      .should('be.visible')
      .and('contain', 'Select a different date:')
  })

  it('should have proper footer content', () => {
    cy.visit(uiServiceUrl)
    
    cy.get('[data-testid="footer-text"]')
      .should('be.visible')
      .and('contain', 'Powered by inspiration • Refreshed daily')
  })

  it('should verify quote service API is accessible', () => {
    // Test the quote service API directly
    cy.request('GET', `${quoteServiceUrl}/api/quotes/today`)
      .then((response) => {
        expect(response.status).to.eq(200)
        expect(response.body).to.have.property('text')
        expect(response.body).to.have.property('author')
        expect(response.body.text).to.be.a('string')
        expect(response.body.author).to.be.a('string')
      })
  })

  it('should verify UI service is accessible', () => {
    // Test the UI service directly
    cy.request('GET', uiServiceUrl)
      .then((response) => {
        expect(response.status).to.eq(200)
        expect(response.body).to.contain('Quotopia')
        expect(response.body).to.contain('data-testid="quotopia-container"')
      })
  })
})
