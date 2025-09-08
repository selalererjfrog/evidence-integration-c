package com.example.quotefday.service;

import com.example.quotefday.model.Quote;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.time.LocalDate;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class QuoteServiceTest {

    private QuoteService quoteService;

    @BeforeEach
    void setUp() {
        quoteService = new QuoteService();
    }

    @Test
    void getQuoteOfTheDay_ShouldReturnQuote() {
        // When
        Quote quote = quoteService.getQuoteOfTheDay();

        // Then
        assertNotNull(quote);
        assertNotNull(quote.getText());
        assertNotNull(quote.getAuthor());
        assertNotNull(quote.getDate());
        assertEquals(LocalDate.now(), quote.getDate());
    }

    @Test
    void getQuoteOfTheDay_WithSpecificDate_ShouldReturnQuoteForThatDate() {
        // Given
        LocalDate testDate = LocalDate.of(2024, 1, 15);

        // When
        Quote quote = quoteService.getQuoteOfTheDay(testDate);

        // Then
        assertNotNull(quote);
        assertNotNull(quote.getText());
        assertNotNull(quote.getAuthor());
        assertEquals(testDate, quote.getDate());
    }

    @Test
    void getQuoteOfTheDay_SameDate_ShouldReturnSameQuote() {
        // Given
        LocalDate testDate = LocalDate.of(2024, 1, 15);

        // When
        Quote quote1 = quoteService.getQuoteOfTheDay(testDate);
        Quote quote2 = quoteService.getQuoteOfTheDay(testDate);

        // Then
        assertEquals(quote1.getText(), quote2.getText());
        assertEquals(quote1.getAuthor(), quote2.getAuthor());
    }

    @Test
    void getAllQuotes_ShouldReturnAllQuotes() {
        // When
        List<Quote> quotes = quoteService.getAllQuotes();

        // Then
        assertNotNull(quotes);
        assertFalse(quotes.isEmpty());
        assertTrue(quotes.size() > 0);
        
        // Verify each quote has required fields
        quotes.forEach(quote -> {
            assertNotNull(quote.getText());
            assertNotNull(quote.getAuthor());
            assertNotNull(quote.getDate());
        });
    }

    @Test
    void getAllQuotes_ShouldReturnQuotesWithTodayDate() {
        // When
        List<Quote> quotes = quoteService.getAllQuotes();

        // Then
        LocalDate today = LocalDate.now();
        quotes.forEach(quote -> assertEquals(today, quote.getDate()));
    }
}
