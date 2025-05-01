package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"blog/internal/service"
	"blog/templates"
)

// Handler contains all HTTP handlers and their dependencies
type Handler struct {
	PostService *service.PostService
	BaseURL     string
	SiteName    string
}

// NewHandler creates a new handler
func NewHandler(postService *service.PostService, baseURL, siteName string) *Handler {
	return &Handler{
		PostService: postService,
		BaseURL:     baseURL,
		SiteName:    siteName,
	}
}

// Register registers all routes with the provided router
func (h *Handler) Register(r chi.Router) {
	r.Get("/", h.HomeHandler)
	r.Get("/me", h.MeHandler)
	r.Get("/portfolio", h.PortfolioHandler)
	r.Get("/books", h.BooksHandler)
	r.Get("/languages", h.LanguagesHandler)
	r.Get("/tools", h.ToolsHandler)
	r.Get("/posts/{slug}", h.PostHandler)
	r.Get("/rss", h.RSSHandler)
}

// HomeHandler handles the home page
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	posts := h.PostService.GetRecentPosts(10)
	renderView(w, templates.HomePage(h.SiteName, posts))
}

// MeHandler handles the about me page
func (h *Handler) MeHandler(w http.ResponseWriter, r *http.Request) {
	renderView(w, templates.MePage(h.SiteName))
}

// PortfolioHandler handles the portfolio page
func (h *Handler) PortfolioHandler(w http.ResponseWriter, r *http.Request) {
	renderView(w, templates.PortfolioPage(h.SiteName))
}

// BooksHandler handles the books page
func (h *Handler) BooksHandler(w http.ResponseWriter, r *http.Request) {
	posts := h.PostService.GetPostsByTag("book")
	renderView(w, templates.BooksPage(h.SiteName, posts))
}

// LanguagesHandler handles the languages page
func (h *Handler) LanguagesHandler(w http.ResponseWriter, r *http.Request) {
	renderView(w, templates.LanguagesPage(h.SiteName))
}

// ToolsHandler handles the developer tools page
func (h *Handler) ToolsHandler(w http.ResponseWriter, r *http.Request) {
	renderView(w, templates.ToolsPage(h.SiteName))
}

// PostHandler handles single post pages
func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	post, err := h.PostService.GetPostBySlug(slug)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	renderView(w, templates.PostPage(h.SiteName, post))
}

// RSSHandler handles the RSS feed
func (h *Handler) RSSHandler(w http.ResponseWriter, r *http.Request) {
	rss := h.PostService.GenerateRSS(h.SiteName, "Latest blog posts", h.BaseURL)
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(rss))
}

// Helper function to render templ views
func renderView(w http.ResponseWriter, component templ.Component) {
	err := component.Render(context.Background(), w)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
} 