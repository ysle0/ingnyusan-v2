package middleware

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ingnyusan/v2/internal/cache"
)

// CachedResponseWriter is a custom ResponseWriter that captures the response
type CachedResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

// Write captures the response and writes it to the underlying buffer
func (c *CachedResponseWriter) Write(b []byte) (int, error) {
	c.buf.Write(b)
	return c.ResponseWriter.Write(b)
}

// CachePage caches entire page responses with the given TTL
func CachePage(c cache.Cache, ttl time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip caching for non-GET requests
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			// Skip caching for API requests
			if strings.HasPrefix(r.URL.Path, "/api/") {
				next.ServeHTTP(w, r)
				return
			}

			// Skip caching for the editor page
			if r.URL.Path == "/editor" {
				next.ServeHTTP(w, r)
				return
			}

			// Create cache key from the request URL
			cacheKey := fmt.Sprintf("page:%s", r.URL.String())

			// Try to get the cached response
			cachedResponse, err := c.Get(context.Background(), cacheKey)
			if err == nil {
				// Cache hit - return the cached response
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Header().Set("X-Cache", "HIT")
				w.Write([]byte(cachedResponse))
				return
			}

			// Cache miss - capture the response
			buf := &bytes.Buffer{}
			cw := &CachedResponseWriter{
				ResponseWriter: w,
				buf:            buf,
			}

			next.ServeHTTP(cw, r)

			// Only cache successful responses with HTML content
			contentType := w.Header().Get("Content-Type")
			if strings.Contains(contentType, "text/html") && buf.Len() > 0 {
				// Store the response in the cache
				c.Set(context.Background(), cacheKey, buf.String(), ttl)
			}
		})
	}
}

// CacheInvalidator returns a middleware that invalidates cache entries
func CacheInvalidator(c cache.Cache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Execute the next handler first
			next.ServeHTTP(w, r)

			// If this is a POST to save a post, invalidate caches
			if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/save-post") {
				// Invalidate home page cache
				c.Delete(context.Background(), "page:/")

				// Invalidate the post page if we can extract the slug
				_ = r.ParseForm()
				if slug := r.FormValue("slug"); slug != "" {
					c.Delete(context.Background(), fmt.Sprintf("page:/posts/%s", slug))
				}

				// Invalidate any tag-based pages that might be affected
				if tags := r.FormValue("tags"); tags != "" {
					for tag := range strings.SplitSeq(tags, ",") {
						tag = strings.TrimSpace(tag)
						if tag == "book" {
							c.Delete(context.Background(), "page:/books")
						}
					}
				}
			}
		})
	}
}
