package com.example.quotefday.model;

import io.swagger.v3.oas.annotations.media.Schema;
import java.time.LocalDate;

@Schema(description = "Quote entity representing an inspirational quote")
public class Quote {
    @Schema(description = "The inspirational quote text", example = "The only person you are destined to become is the person you decide to be.")
    private String text;
    
    @Schema(description = "The author of the quote", example = "Ralph Waldo Emerson")
    private String author;
    
    @Schema(description = "The date associated with this quote", example = "2025-08-10")
    private LocalDate date;

    public Quote() {
    }

    public Quote(String text, String author, LocalDate date) {
        this.text = text;
        this.author = author;
        this.date = date;
    }

    public String getText() {
        return text;
    }

    public void setText(String text) {
        this.text = text;
    }

    public String getAuthor() {
        return author;
    }

    public void setAuthor(String author) {
        this.author = author;
    }

    public LocalDate getDate() {
        return date;
    }

    public void setDate(LocalDate date) {
        this.date = date;
    }

    @Override
    public String toString() {
        return "Quote{" +
                "text='" + text + '\'' +
                ", author='" + author + '\'' +
                ", date=" + date +
                '}';
    }
}
