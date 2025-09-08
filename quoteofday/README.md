# Quote of Day Service

A Spring Boot REST API service that provides inspirational quotes of the day. The service returns different quotes based on the current date, ensuring users get a new quote each day.

## Features

- **Quote of the Day**: Get today's inspirational quote
- **Historical Quotes**: Get quotes for specific dates
- **All Quotes**: Retrieve all available quotes
- **Health Check**: Service health endpoint
- **RESTful API**: Clean REST endpoints with JSON responses

## Technology Stack

- **Java 17**
- **Spring Boot 3.2.0**
- **Maven** (Build tool)
- **JUnit 5** (Testing framework)
- **Spring Boot Test** (Integration testing)

## Project Structure

```
src/
├── main/
│   ├── java/com/example/quotefday/
│   │   ├── QuoteOfDayApplication.java    # Main application class
│   │   ├── controller/
│   │   │   └── QuoteController.java      # REST API endpoints
│   │   ├── service/
│   │   │   └── QuoteService.java         # Business logic
│   │   └── model/
│   │       └── Quote.java                # Data model
│   └── resources/
│       └── application.properties        # Configuration
└── test/
    └── java/com/example/quotefday/
        ├── QuoteOfDayApplicationTests.java    # Application tests
        ├── controller/
        │   └── QuoteControllerTest.java       # Controller tests
        └── service/
            └── QuoteServiceTest.java          # Service tests
```

## Getting Started

### Prerequisites

- Java 17 or higher
- Maven 3.6 or higher

### Building the Project

```bash
# Navigate to the service directory
cd quoteofday

# Build the project
mvn clean install
```

### Running the Application

```bash
# Run the application
mvn spring-boot:run
```

The application will start on `http://localhost:8001`

### Running Tests

```bash
# Run all tests
mvn test

# Run tests with coverage
mvn test jacoco:report
```

## API Endpoints

### 1. Get Today's Quote

**GET** `/api/quotes/today`

Returns the quote of the day for the current date.

**Response:**
```json
{
  "text": "The only way to do great work is to love what you do.",
  "author": "Steve Jobs",
  "date": "2024-01-15"
}
```

### 2. Get Quote for Specific Date

**GET** `/api/quotes/date/{date}`

Returns the quote for a specific date (format: YYYY-MM-DD).

**Example:** `GET /api/quotes/date/2024-01-15`

**Response:**
```json
{
  "text": "Life is what happens when you're busy making other plans.",
  "author": "John Lennon",
  "date": "2024-01-15"
}
```

### 3. Get All Quotes

**GET** `/api/quotes`

Returns all available quotes with today's date.

**Response:**
```json
[
  {
    "text": "The only way to do great work is to love what you do.",
    "author": "Steve Jobs",
    "date": "2024-01-15"
  },
  {
    "text": "Life is what happens when you're busy making other plans.",
    "author": "John Lennon",
    "date": "2024-01-15"
  }
]
```

### 4. Health Check

**GET** `/api/quotes/health`

Returns a simple health check message.

**Response:**
```
Quote of Day Service is running!
```

## Testing the API

### Using curl

```bash
# Get today's quote
curl http://localhost:8001/api/quotes/today

# Get quote for specific date
curl http://localhost:8001/api/quotes/date/2024-01-15

# Get all quotes
curl http://localhost:8001/api/quotes

# Health check
curl http://localhost:8001/api/quotes/health
```

### Using a web browser

Simply navigate to:
- `http://localhost:8001/api/quotes/today`
- `http://localhost:8001/api/quotes/date/2024-01-15`
- `http://localhost:8001/api/quotes`
- `http://localhost:8001/api/quotes/health`

## How It Works

The service uses a deterministic algorithm to select quotes based on the day of the year:

1. **Quote Selection**: The service calculates the day of the year (1-366) and uses modulo operation to select a quote from the predefined list
2. **Consistency**: The same date will always return the same quote, ensuring consistency
3. **Daily Rotation**: Different dates return different quotes, providing variety

## Configuration

The application can be configured through `src/main/resources/application.properties`:

- **Server Port**: Default is 8001
- **Logging**: Configured for INFO level
- **Jackson**: Configured to exclude null values and format dates properly

## Development

### Adding New Quotes

To add new quotes, modify the `quotes` list in `QuoteService.java`:

```java
private final List<Quote> quotes = Arrays.asList(
    new Quote("Your new quote here", "Author Name", LocalDate.now()),
    // ... existing quotes
);
```

### Running in Development Mode

The application includes Spring Boot DevTools for enhanced development experience:

- Automatic restart on code changes
- Live reload capabilities
- Enhanced error pages

## Docker Support

The service includes a multi-stage Dockerfile for optimized container builds:

```bash
# Build Docker image
docker build -t quote-of-day-service .

# Run container
docker run -p 8001:8001 quote-of-day-service
```

### Docker Features

✅ **Multi-stage build**: Optimized image size  
✅ **Security**: Non-root user execution  
✅ **Health checks**: Built-in health monitoring  
✅ **Alpine base**: Minimal attack surface  

## Deployment

### Building for Production

```bash
# Create executable JAR
mvn clean package

# The JAR file will be created in target/ directory
java -jar target/quote-of-day-service-1.0.0.jar
```

### Container Deployment

```bash
# Build and run with Docker
docker build -t quote-of-day-service .
docker run -p 8001:8001 quote-of-day-service
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the MIT License.
