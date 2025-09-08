package com.example.quotefday.controller;

import com.example.quotefday.model.Quote;
import com.example.quotefday.service.QuoteService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.ExampleObject;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDate;
import java.util.List;

@RestController
@RequestMapping("/api/quotes")
@CrossOrigin(origins = "*")
@Tag(name = "Quote of Day", description = "Quote of Day management APIs")
public class QuoteController {

    private final QuoteService quoteService;

    @Autowired
    public QuoteController(QuoteService quoteService) {
        this.quoteService = quoteService;
    }

    @GetMapping("/today")
    @Operation(
        summary = "Get today's quote",
        description = "Retrieves the inspirational quote of the day for the current date"
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved today's quote",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = Quote.class),
                examples = @ExampleObject(
                    name = "Sample Quote",
                    value = "{\"text\":\"The only person you are destined to become is the person you decide to be.\",\"author\":\"Ralph Waldo Emerson\",\"date\":\"2025-08-10\"}"
                )
            )
        ),
        @ApiResponse(
            responseCode = "500",
            description = "Internal server error",
            content = @Content
        )
    })
    public ResponseEntity<Quote> getQuoteOfTheDay() {
        Quote quote = quoteService.getQuoteOfTheDay();
        return ResponseEntity.ok(quote);
    }

    @GetMapping("/date/{date}")
    @Operation(
        summary = "Get quote by date",
        description = "Retrieves the inspirational quote for a specific date"
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved quote for the specified date",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = Quote.class)
            )
        ),
        @ApiResponse(
            responseCode = "400",
            description = "Invalid date format",
            content = @Content
        ),
        @ApiResponse(
            responseCode = "404",
            description = "Quote not found for the specified date",
            content = @Content
        )
    })
    public ResponseEntity<Quote> getQuoteOfTheDay(
            @Parameter(description = "Date in ISO format (YYYY-MM-DD)", example = "2025-08-10")
            @PathVariable @DateTimeFormat(iso = DateTimeFormat.ISO.DATE) LocalDate date) {
        Quote quote = quoteService.getQuoteOfTheDay(date);
        return ResponseEntity.ok(quote);
    }

    @GetMapping
    @Operation(
        summary = "Get all quotes",
        description = "Retrieves all available inspirational quotes in the system"
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Successfully retrieved all quotes",
            content = @Content(
                mediaType = "application/json",
                schema = @Schema(implementation = Quote.class)
            )
        ),
        @ApiResponse(
            responseCode = "500",
            description = "Internal server error",
            content = @Content
        )
    })
    public ResponseEntity<List<Quote>> getAllQuotes() {
        List<Quote> quotes = quoteService.getAllQuotes();
        return ResponseEntity.ok(quotes);
    }

    @GetMapping("/health")
    @Operation(
        summary = "Health check",
        description = "Checks if the Quote of Day service is running and healthy"
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200",
            description = "Service is healthy",
            content = @Content(
                mediaType = "text/plain",
                examples = @ExampleObject(
                    name = "Health Response",
                    value = "Quote of Day Service is running!"
                )
            )
        )
    })
    public ResponseEntity<String> health() {
        return ResponseEntity.ok("Quote of Day Service is running!");
    }
}
