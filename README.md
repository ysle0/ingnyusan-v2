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

## Configuration

The application uses a TOML configuration file by default. Create a `config.toml` file in the project root with the following structure:

```toml
# Server settings
port = "8080"

# Content settings
posts_dir = "./posts"
site_name = "My Blog"

# URLs and endpoints
base_url = "http://localhost:8080"

# Redis cache settings
redis_url = "localhost:6379"
cache_ttl = "10m" # Duration format: 10m = 10 minutes
```

You can specify a different configuration file path using the `-config` flag:

```
go run cmd/server/main.go -config path/to/config.toml
```

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
language: english
---
```

### Time-to-Read Calculation

The blog automatically calculates the estimated reading time based on the content:

- For English content (default): Uses a reading speed of 200 words per minute
- For Korean content: Uses a reading speed of 180 characters per minute (set `language: korean` in frontmatter)

To specify the language of your post, add the `language` field to the frontmatter:

```
language: korean  # For Korean content
language: english # For English content (default if not specified)
```

## License

MIT 