package com.example.quotefday.controller;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.boot.test.web.server.LocalServerPort;
import org.springframework.http.ResponseEntity;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class QuoteControllerTest {

    @LocalServerPort
    private int port;

    @Autowired
    private TestRestTemplate restTemplate;

    @Test
    void getQuoteOfTheDay_ShouldReturnQuote() {
        // When
        ResponseEntity<String> response = restTemplate.getForEntity(
                "http://localhost:" + port + "/api/quotes/today", String.class);

        // Then
        assertEquals(200, response.getStatusCodeValue());
        assertTrue(response.getBody().contains("text"));
        assertTrue(response.getBody().contains("author"));
        assertTrue(response.getBody().contains("date"));
    }

    @Test
    void getQuoteOfTheDay_WithDate_ShouldReturnQuoteForDate() {
        // When
        ResponseEntity<String> response = restTemplate.getForEntity(
                "http://localhost:" + port + "/api/quotes/date/2024-01-15", String.class);

        // Then
        assertEquals(200, response.getStatusCodeValue());
        assertTrue(response.getBody().contains("text"));
        assertTrue(response.getBody().contains("author"));
        assertTrue(response.getBody().contains("2024-01-15"));
    }

    @Test
    void getAllQuotes_ShouldReturnAllQuotes() {
        // When
        ResponseEntity<String> response = restTemplate.getForEntity(
                "http://localhost:" + port + "/api/quotes", String.class);

        // Then
        assertEquals(200, response.getStatusCodeValue());
        assertTrue(response.getBody().startsWith("["));
        assertTrue(response.getBody().endsWith("]"));
    }

    @Test
    void health_ShouldReturnHealthMessage() {
        // When
        ResponseEntity<String> response = restTemplate.getForEntity(
                "http://localhost:" + port + "/api/quotes/health", String.class);

        // Then
        assertEquals(200, response.getStatusCodeValue());
        assertEquals("Quote of Day Service is running!", response.getBody());
    }
}
