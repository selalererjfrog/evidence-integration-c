package com.example.quotefday.service;

import com.example.quotefday.model.Quote;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.util.Arrays;
import java.util.List;

@Service
public class QuoteService {

    private final List<Quote> quotes = Arrays.asList(
            new Quote("The only way to do great work is to love whats you do.", "Steve Jobs", LocalDate.now()),
            new Quote("Life is what happens when you're busy making other plans.", "John Lennon", LocalDate.now()),
            new Quote("The future belongs to those who believe in the beauty of their dreams.", "Eleanor Roosevelt", LocalDate.now()),
            new Quote("Success is not final, failure is not fatal: it is the courage to continue that counts.", "Winston Churchill", LocalDate.now()),
            new Quote("The only limit to our realization of tomorrow is our doubts of today.", "Franklin D. Roosevelt", LocalDate.now()),
            new Quote("It does not matter how slowly you go as long as you do not stop.", "Confucius", LocalDate.now()),
            new Quote("The journey of a thousand miles begins with one step.", "Lao Tzu", LocalDate.now()),
            new Quote("What you get by achieving your goals is not as important as what you become by achieving your goals.", "Zig Ziglar", LocalDate.now()),
            new Quote("The best way to predict the future is to invent it.", "Alan Kay", LocalDate.now()),
            new Quote("Don't watch the clock; do what it does. Keep going.", "Sam Levenson", LocalDate.now()),
            new Quote("Believe you can and you're halfway there.", "Theodore Roosevelt", LocalDate.now()),
            new Quote("The mind is everything. What you think you become.", "Buddha", LocalDate.now()),
            new Quote("The only person you are destined to become is the person you decide to be.", "Ralph Waldo Emerson", LocalDate.now()),
            new Quote("Everything you've ever wanted is on the other side of fear.", "George Addair", LocalDate.now()),
            new Quote("The way to get started is to quit talking and begin doing.", "Walt Disney", LocalDate.now()),
            new Quote("Your time is limited, don't waste it living someone else's life.", "Steve Jobs", LocalDate.now()),
            new Quote("The greatest glory in living lies not in never falling, but in rising every time we fall.", "Nelson Mandela", LocalDate.now()),
            new Quote("In the middle of difficulty lies opportunity.", "Albert Einstein", LocalDate.now()),
            new Quote("The only impossible journey is the one you never begin.", "Tony Robbins", LocalDate.now()),
            new Quote("What you do today can improve all your tomorrows.", "Ralph Marston", LocalDate.now()),
            new Quote("The secret of getting ahead is getting started.", "Mark Twain", LocalDate.now()),
            new Quote("Don't let yesterday take up too much of today.", "Will Rogers", LocalDate.now()),
            new Quote("The harder you work for something, the greater you'll feel when you achieve it.", "Unknown", LocalDate.now()),
            new Quote("Dream big and dare to fail.", "Norman Vaughan", LocalDate.now()),
            new Quote("The best revenge is massive success.", "Frank Sinatra", LocalDate.now()),
            new Quote("I find that the harder I work, the more luck I seem to have.", "Thomas Jefferson", LocalDate.now()),
            new Quote("Success is walking from failure to failure with no loss of enthusiasm.", "Winston Churchill", LocalDate.now()),
            new Quote("The only place where success comes before work is in the dictionary.", "Vidal Sassoon", LocalDate.now()),
            new Quote("Don't be afraid to give up the good to go for the great.", "John D. Rockefeller", LocalDate.now()),
            new Quote("The future depends on what you do today.", "Mahatma Gandhi", LocalDate.now())
    );

    public Quote getQuoteOfTheDay() {
        int dayOfYear = LocalDate.now().getDayOfYear();
        int index = dayOfYear % quotes.size();
        Quote quote = quotes.get(index);
        quote.setDate(LocalDate.now());
        return quote;
    }

    public Quote getQuoteOfTheDay(LocalDate date) {
        int dayOfYear = date.getDayOfYear();
        int index = dayOfYear % quotes.size();
        Quote quote = quotes.get(index);
        quote.setDate(date);
        return quote;
    }

    public List<Quote> getAllQuotes() {
        return quotes.stream()
                .map(quote -> new Quote(quote.getText(), quote.getAuthor(), LocalDate.now()))
                .toList();
    }
}
