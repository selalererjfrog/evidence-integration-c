const { defineConfig } = require('cypress')
const { exec } = require('child_process')
const { promisify } = require('util')

const execAsync = promisify(exec)

// Helper function to check service health
async function checkServiceHealth(url, endpoint = '/actuator/health') {
  try {
    const { stdout } = await execAsync(`curl -f -s ${url}${endpoint}`)
    return { running: true, response: stdout }
  } catch (error) {
    return { running: false, error: error.message }
  }
}

module.exports = defineConfig({
  e2e: {
    baseUrl: 'http://localhost:8001',
    supportFile: 'cypress/support/e2e.js',
    specPattern: 'cypress/e2e/**/*.cy.js',
    video: false,
    screenshotOnRunFailure: false,
    env: {
      quoteServiceUrl: 'http://localhost:8001',
      uiServiceUrl: 'http://localhost:8081'
    },
    // Performance optimizations
    defaultCommandTimeout: 10000,
    requestTimeout: 10000,
    responseTimeout: 10000,
    pageLoadTimeout: 30000,
    
    // Reporter configuration for JSON output
    reporter: 'json',
    reporterOptions: {
      outputFile: 'cypress/results/results.json'
    },
    
    setupNodeEvents(on, config) {
      // Task to check if service is running
      on('task', {
        async checkServiceHealth(url, endpoint = '/actuator/health') {
          return await checkServiceHealth(url, endpoint)
        },
        
        async startServicesIfNeeded({ quoteServiceUrl, uiServiceUrl }) {
          console.log('🔍 Checking if services are running...')
          
          // Check quote service
          const quoteHealth = await checkServiceHealth(quoteServiceUrl)
          if (!quoteHealth.running) {
            console.log('🚀 Starting quoteofday service...')
            try {
              await execAsync('cd ../quoteofday && ./mvnw spring-boot:run -Dspring-boot.run.profiles=qa', { 
                cwd: process.cwd(),
                stdio: 'pipe'
              })
              console.log('✅ Quoteofday service started')
            } catch (error) {
              console.log('⚠️ Quoteofday service may already be running or failed to start')
            }
          } else {
            console.log('✅ Quoteofday service is already running')
          }
          
          // Check UI service (we'll assume it's running via Docker or manually)
          const uiHealth = await checkServiceHealth(uiServiceUrl, '/')
          if (!uiHealth.running) {
            console.log('⚠️ UI service is not running. Please start it manually or via Docker.')
            console.log('   You can start it with: cd ../quotopia-ui && docker run -p 8081:80 quotopia-ui:latest')
          } else {
            console.log('✅ UI service is running')
          }
          
          return { quoteService: quoteHealth.running, uiService: uiHealth.running }
        },
        
        async stopService(serviceName) {
          console.log(`🛑 Stopping ${serviceName} service...`)
          try {
            if (serviceName === 'quoteofday') {
              // Find and kill the Java process
              await execAsync("pkill -f 'spring-boot:run'")
              console.log('✅ Quoteofday service stopped')
            }
            return { success: true }
          } catch (error) {
            console.log('⚠️ Service may not be running or already stopped')
            return { success: false, error: error.message }
          }
        },
        
        async startService(serviceName) {
          console.log(`🚀 Starting ${serviceName} service...`)
          try {
            if (serviceName === 'quoteofday') {
              await execAsync('cd ../quoteofday && ./mvnw spring-boot:run -Dspring-boot.run.profiles=qa', {
                cwd: process.cwd(),
                stdio: 'pipe'
              })
              console.log('✅ Quoteofday service started')
            }
            return { success: true }
          } catch (error) {
            console.log('⚠️ Service may already be running or failed to start')
            return { success: false, error: error.message }
          }
        }
      })
    }
  }
})
