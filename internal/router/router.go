package router

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"ingnyusan/v2/internal/cache"
	"ingnyusan/v2/internal/config"
	"ingnyusan/v2/internal/handler"
	"ingnyusan/v2/internal/i18n"
	customMiddleware "ingnyusan/v2/internal/middleware"
	"ingnyusan/v2/internal/service"
	"ingnyusan/v2/templates"
)

// SetupRouter configures the Chi router with all routes and middleware
func SetupRouter(
	cfg *config.Config,
	postSvc *service.PostService,
	cacher *cache.Cacher,
) chi.Router {
	// Set up translator
	translator := i18n.NewTranslator(cfg.I18n.TranslationsDir, cfg.I18n.DefaultLang)

	// Set global config for templates
	templates.SetConfig(cfg)

	h := handler.NewHandler(postSvc, cfg.Server.BaseURL, cfg.Content.SiteName, translator, cfg)
	r := chi.NewRouter()

	setupStdMiddleware(r)

	// Add i18n middleware
	r.Use(customMiddleware.I18nMiddleware(translator))

	// Add caching middleware if Redis is available
	if cacher != nil {
		setupCachingMiddleware(r, cacher, cfg.Redis.CacheTTL)
	}

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Register routes
	h.Register(r)
	h.RegisterAPI(r)

	// Add editor route
	r.Get("/editor", func(w http.ResponseWriter, r *http.Request) {
		renderTemplView(w, r, templates.EditorPage(cfg.Content.SiteName, translator, r))
	})

	return r
}

// setupStdMiddleware adds the standard middleware to the router
func setupStdMiddleware(r chi.Router) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

// setupCachingMiddleware adds Redis caching middleware to the router
func setupCachingMiddleware(r chi.Router, cacher *cache.Cacher, cacheTTL time.Duration) {
	// Add cache invalidator first (to run after the handlers)
	r.Use(customMiddleware.CacheInvalidator(cacher))
	// Add page caching middleware
	r.Use(customMiddleware.CachePage(cacher, cacheTTL))
}

// renderTemplView renders a templ component to the response writer
func renderTemplView(w http.ResponseWriter, r *http.Request, component templ.Component) {
	err := component.Render(context.Background(), w)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
