// Quote of the day functionality
class QuotopiaApp {
    constructor() {
        // Use localhost when running locally, host.docker.internal when running in Docker
        this.baseUrl = window.location.hostname === 'localhost' ? 'http://localhost:8001/api/quotes' : 'http://host.docker.internal:8001/api/quotes';
        this.quoteTextElement = document.getElementById('quoteText');
        this.quoteAuthorElement = document.getElementById('quoteAuthor');
        this.quoteDateElement = document.getElementById('quoteDate');
        this.dateSelector = document.getElementById('dateSelector');
        
        this.init();
    }

    async init() {
        try {
            await this.loadQuote();
            this.setupDateSelector();
        } catch (error) {
            console.error('Error loading quote:', error);
            this.showError();
        }
    }

    async loadQuote(date = null) {
        try {
            const url = date ? `${this.baseUrl}/date/${date}` : `${this.baseUrl}/today`;
            console.log('Fetching quote from:', url);
            console.log('Base URL:', this.baseUrl);
            
            const response = await fetch(url);
            console.log('Response status:', response.status);
            console.log('Response headers:', response.headers);
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const data = await response.json();
            console.log('Quote data received:', data);
            this.displayQuote(data);
        } catch (error) {
            console.error('Failed to fetch quote:', error);
            console.error('Error details:', error.message);
            this.showError();
        }
    }

    setupDateSelector() {
        if (this.dateSelector) {
            this.dateSelector.addEventListener('change', async (event) => {
                const selectedDate = event.target.value;
                if (selectedDate) {
                    await this.loadQuote(selectedDate);
                }
            });
        }
    }

    displayQuote(quoteData) {
        // Add fade-in animation
        this.quoteTextElement.style.opacity = '0';
        this.quoteAuthorElement.style.opacity = '0';
        this.quoteDateElement.style.opacity = '0';

        // Update content
        this.quoteTextElement.textContent = quoteData.text;
        this.quoteAuthorElement.textContent = `— ${quoteData.author}`;
        
        // Format and display date
        if (quoteData.date) {
            const date = new Date(quoteData.date);
            const formattedDate = date.toLocaleDateString('en-US', {
                weekday: 'long',
                year: 'numeric',
                month: 'long',
                day: 'numeric'
            });
            this.quoteDateElement.textContent = formattedDate;
        }

        // Fade in the content
        setTimeout(() => {
            this.quoteTextElement.style.transition = 'opacity 0.8s ease-in-out';
            this.quoteAuthorElement.style.transition = 'opacity 0.8s ease-in-out';
            this.quoteDateElement.style.transition = 'opacity 0.8s ease-in-out';
            
            this.quoteTextElement.style.opacity = '1';
            this.quoteAuthorElement.style.opacity = '1';
            this.quoteDateElement.style.opacity = '1';
        }, 100);
    }

    showError() {
        this.quoteTextElement.innerHTML = `
            <div style="display: flex; align-items: center; justify-content: center; gap: 10px;">
                <div class="loading">
                    <div></div>
                    <div></div>
                </div>
                <span>Unable to load today's inspiration</span>
            </div>
        `;
        this.quoteAuthorElement.textContent = 'Please try again later';
        this.quoteDateElement.textContent = '';
    }

    // Method to refresh quote (can be called manually if needed)
    async refreshQuote() {
        this.quoteTextElement.textContent = 'Loading today\'s inspiration...';
        this.quoteAuthorElement.textContent = '—';
        this.quoteDateElement.textContent = '';
        
        await this.loadQuote();
    }

    // Method to load quote for specific date
    async loadQuoteForDate(date) {
        this.quoteTextElement.textContent = 'Loading inspiration...';
        this.quoteAuthorElement.textContent = '—';
        this.quoteDateElement.textContent = '';
        
        await this.loadQuote(date);
    }
}

// Initialize the app when DOM is loaded
let quotopiaApp;

document.addEventListener('DOMContentLoaded', () => {
    quotopiaApp = new QuotopiaApp();
    
    const quoteCard = document.querySelector('.quote-card');
    
    // Add click to refresh functionality
    quoteCard.addEventListener('click', async () => {
        // Add a subtle click effect
        quoteCard.style.transform = 'scale(0.98)';
        setTimeout(() => {
            quoteCard.style.transform = '';
        }, 150);
        
        // Refresh the quote using the existing app instance
        await quotopiaApp.refreshQuote();
    });
    
    // Add hover effect for better UX
    quoteCard.style.cursor = 'pointer';
    
    // Add a subtle tooltip
    quoteCard.title = 'Click to refresh quote';
});
