# Blog Website Plan

## Project Structure
```
├── cmd/
│   └── server/
│       └── main.go         # Entry point
├── internal/
│   ├── handler/            # HTTP handlers
│   ├── middleware/         # Middleware functions
│   ├── model/              # Data models
│   ├── parser/             # Markdown parser
│   ├── cache/              # Redis cache implementation
│   └── service/            # Business logic
├── static/                 # Static assets
│   ├── css/
│   └── js/
├── templates/              # templ templates
├── posts/                  # Markdown blog posts
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

## Core Components

1. **HTTP Server (Chi Router)**
   - Main server setup with routing
   - Static file serving
   - CORS configuration
   - Middleware setup

2. **Markdown Parser**
   - Parse .md files to HTML
   - Extract frontmatter (title, date, tags)
   - Support for code highlighting

3. **Blog Post Service**
   - Load posts from embedded filesystem
   - Sort by date
   - Filter by tags/categories
   - Pagination support

4. **Redis Caching**
   - Cache rendered HTML pages
   - Cache individual blog posts
   - Cache API responses
   - Configurable TTL (time-to-live)
   - Cache invalidation on post updates

5. **In-Editor Post Writer**
   - Simple web interface for writing posts
   - Live preview with markdown rendering
   - Save directly to filesystem

6. **Time-to-Read Calculator**
   - Estimate reading time based on word count
   - Display on post pages

7. **RSS Feed Generator**
   - Generate RSS XML
   - Update on new posts

## Implementation Steps

1. Set up project structure and dependencies
2. Create basic Chi router with routes
3. Build markdown parser and post service
4. Implement Redis caching for rendered HTML
5. Design and implement templates with templ
6. Create minimal CSS with chosen framework
7. Implement editor interface with htmx
8. Add RSS feed generator
9. Set up Docker configuration with Redis service
10. Integrate Google Analytics
11. Deploy and test

## Tech Specifics

- **Go-Chi**: For routing and handling HTTP requests
- **htmx**: For dynamic content without JavaScript
- **templ**: For type-safe HTML templates
- **CSS Framework**: [Consider: Water.css, Pico.css, or Bootstrap - minimalist options]
- **Redis**: For caching rendered HTML and improving performance
- **Docker/Docker-Compose**: For containerization
- **Google Analytics**: For visitor tracking 