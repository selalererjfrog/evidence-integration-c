# Quotopia - Inspire yourself THEN inspire the world

A beautiful, inspiring static web page that displays the quote of the day from your microservice API.

## Features

- ✨ Beautiful, modern design with gradient backgrounds
- 💡 Animated header with inspiration-themed illustrations
- 📱 Fully responsive design for all devices
- 🔄 Interactive quote card (click to refresh)
- 🌟 Smooth animations and transitions
- 📅 Displays formatted date information

## Files

- `index.html` - Main HTML structure
- `styles.css` - Beautiful CSS styling with animations
- `script.js` - JavaScript for API integration and interactivity

## Setup

1. **Start your quote microservice** on `http://localhost:8080`
2. **Open `index.html`** in your web browser
3. **Enjoy your daily inspiration!**

## API Requirements

The web page expects your microservice to provide quotes at:
```
GET http://localhost:8080/api/quotes/today
```

Expected response format:
```json
{
  "text": "The only way to do great work is to love what you do.",
  "author": "Steve Jobs",
  "date": "2025-08-28"
}
```

## Features

### Design Elements
- **Header**: Quotopia title with subtitle "Inspire yourself THEN inspire the world"
- **Header Image**: Animated light bulb with sparkles representing inspiration
- **Quote Card**: Clean, modern card displaying the quote, author, and date
- **Footer**: Simple footer with branding

### Interactivity
- **Click to Refresh**: Click anywhere on the quote card to fetch a new quote
- **Hover Effects**: Subtle animations on hover
- **Loading States**: Beautiful loading animations when fetching quotes
- **Error Handling**: Graceful error display if API is unavailable

### Responsive Design
- **Mobile-friendly**: Optimized for phones and tablets
- **Desktop-optimized**: Beautiful layout on larger screens
- **Flexible**: Adapts to different screen sizes

## Customization

You can easily customize the design by modifying:
- **Colors**: Update the gradient colors in `styles.css`
- **Fonts**: Change the Google Fonts in `index.html`
- **API URL**: Modify the `apiUrl` in `script.js`
- **Animations**: Adjust timing and effects in `styles.css`

## Browser Compatibility

- ✅ Chrome (recommended)
- ✅ Firefox
- ✅ Safari
- ✅ Edge

## Troubleshooting

If the quote doesn't load:
1. Ensure your microservice is running on `http://localhost:8080`
2. Check that the API endpoint `/api/quotes/today` is accessible
3. Verify the response format matches the expected JSON structure
4. Check browser console for any error messages

## License

This is a static web page created for the Quotopia service. Feel free to modify and use as needed.
