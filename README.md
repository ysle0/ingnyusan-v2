# Minimalist Blog

A minimalist blog website built with Go, Chi, HTMX, Templ, and Water.css.

## Features

- Markdown blog posts with frontmatter
- Responsive design with minimal CSS
- Time-to-read calculation
- RSS feed
- In-browser post editor
- Support for tags and categories
- About me page with GitHub integration
- Developer tools section

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

## Running Locally

### Prerequisites

- Go 1.18+
- templ

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/blog.git
   cd blog
   ```

2. Install dependencies:
   ```
   go mod download
   go install github.com/a-h/templ/cmd/templ@latest
   ```

3. Generate templ templates:
   ```
   templ generate
   ```

4. Run the application:
   ```
   go run cmd/server/main.go
   ```

5. Open http://localhost:8080 in your browser

## Docker

You can also run the application using Docker:

```
docker-compose up -d
```

## Adding Blog Posts

1. Create a new markdown file in the `posts` directory
2. Add frontmatter with title, description, date, and tags
3. Write your content in markdown
4. Restart the server or use the in-browser editor at /editor

Example frontmatter:

```
---
title: Hello, World!
description: My first blog post
date: 2023-10-20
tags: introduction, blog
---
```

## License

MIT 