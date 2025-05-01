FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk update && apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate templ files
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog ./cmd/server

# Use a small alpine image for the final image
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/blog /app/blog
COPY --from=builder /app/static /app/static
COPY --from=builder /app/posts /app/posts

# Set working directory
WORKDIR /app

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./blog", "--baseurl", "http://localhost:8080", "--sitename", "My Blog"] 