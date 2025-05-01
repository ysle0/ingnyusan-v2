---
title: Introduction to Templ - A Modern HTML Templating Language for Go
description: Learn about Templ, a type-safe and component-based templating language for Go
date: 2023-12-10
tags: programming, golang, web development, templating
language: english
---

# Introduction to Templ - A Modern HTML Templating Language for Go

When building web applications in Go, templating is a crucial part of generating HTML responses. While Go's standard library offers text/template and html/template packages, modern web development often requires more type safety and component-based architecture. This is where Templ comes in.

## What is Templ?

[Templ](https://templ.guide/) is a templating language designed specifically for building HTML user interfaces in Go. Unlike traditional Go templates, Templ:

- Is type-safe, catching errors at compile-time
- Uses a component-based approach similar to React or Svelte
- Generates efficient Go code
- Works seamlessly with Go's HTTP handlers

## Getting Started with Templ

First, you'll need to install Templ:

```bash
go install github.com/a-h/templ/cmd/templ@latest
```

### Basic Example

Create a file named `greeting.templ`:

```go
package main

templ greeting(name string) {
    <div>
        <h1>Hello, { name }!</h1>
        <p>Welcome to Templ</p>
    </div>
}
```

Then, compile it with:

```bash
templ generate
```

This creates a Go file containing the compiled template function.

### Using Templ in a Go Application

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        component := greeting("World")
        component.Render(r.Context(), w)
    })

    fmt.Println("Server running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

## Component Composition

One of Templ's strengths is component composition:

```go
package main

templ header(title string) {
    <header>
        <h1>{ title }</h1>
        <nav>
            <a href="/">Home</a>
            <a href="/about">About</a>
        </nav>
    </header>
}

templ footer() {
    <footer>
        <p>© 2023 My Website</p>
    </footer>
}

templ layout(title string) {
    <html>
        <head>
            <title>{ title }</title>
            <link rel="stylesheet" href="/static/style.css"/>
        </head>
        <body>
            @header(title)
            <main>
                { children... }
            </main>
            @footer()
        </body>
    </html>
}

templ homePage() {
    @layout("Home Page") {
        <h2>Welcome to our website!</h2>
        <p>This is a simple example of component composition in Templ.</p>
    }
}
```

## Control Flow

Templ supports standard control flow constructs:

### Conditionals

```go
templ conditionalExample(isLoggedIn bool) {
    if isLoggedIn {
        <p>Welcome back, user!</p>
    } else {
        <p>Please log in to continue.</p>
    }
}
```

### Loops

```go
templ loopExample(items []string) {
    <ul>
        for _, item := range items {
            <li>{ item }</li>
        }
    </ul>
}
```

## CSS Integration

Templ makes it easy to add class names and styles:

```go
templ styledComponent(theme string) {
    <div class={ "card " + theme }>
        <h2 class="card-title">Title</h2>
        <p class="card-content">Content</p>
    </div>
}
```

## Dynamic Attributes

```go
templ dynamicAttributes(isDisabled bool) {
    <button disabled?={ isDisabled }>
        Submit
    </button>
}
```

## HTMX Integration

Templ works beautifully with HTMX for building modern, dynamic applications:

```go
templ counterButton(count int) {
    <div>
        <span>Count: { fmt.Sprint(count) }</span>
        <button 
            hx-post="/increment" 
            hx-target="closest div" 
            hx-swap="outerHTML">
            Increment
        </button>
    </div>
}
```

## Why Choose Templ?

- **Type Safety**: Compile-time checking prevents runtime errors
- **Component Model**: Encourages reusable, composable UI components
- **Performance**: Fast rendering with minimal overhead
- **Developer Experience**: Clean syntax that feels natural to Go developers
- **Integration**: Works well with Go's HTTP ecosystem

## Conclusion

Templ represents a modern approach to HTML templating in Go, bringing many of the benefits of component-based frameworks like React to the Go ecosystem. If you're building web applications in Go and find the standard templating limiting, Templ offers a compelling alternative that enhances productivity while maintaining Go's performance benefits. 