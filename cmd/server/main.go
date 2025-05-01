package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"blog/internal/cache"
	"blog/internal/handler"
	customMiddleware "blog/internal/middleware"
	"blog/internal/service"
	"blog/templates"
)

func main() {
	port := flag.String("port", "8080", "Server port")
	postsDir := flag.String("posts", "./posts", "Posts directory")
	baseURL := flag.String("baseurl", "http://localhost:8080", "Base URL for the blog")
	siteName := flag.String("sitename", "My Blog", "Name of the blog")
	redisURL := flag.String("redis", "localhost:6379", "Redis URL for caching")
	cacheTTL := flag.Duration("cachettl", 10*time.Minute, "Cache TTL duration")
	flag.Parse()

	// Initialize the cache
	redisCache := cache.NewRedisCache(*redisURL)
	defer redisCache.Close()

	// Test Redis connection
	_, err := redisCache.Get(context.Background(), "test")
	if err != nil && err.Error() != "redis: nil" {
		log.Printf("WARNING: Could not connect to Redis: %v", err)
		log.Printf("Continuing without caching")
		redisCache = nil
	} else {
		log.Printf("Connected to Redis cache at %s", *redisURL)
	}

	// Initialize the post service
	postService := service.NewPostService(*postsDir)
	err = postService.LoadPosts()
	if err != nil {
		log.Fatalf("Failed to load posts: %v", err)
	}

	// Initialize the handler
	h := handler.NewHandler(postService, *baseURL, *siteName)

	// Set up the router
	r := chi.NewRouter()

	// Middleware
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

	// Add caching middleware if Redis is available
	if redisCache != nil {
		// Add cache invalidator first (to run after the handlers)
		r.Use(customMiddleware.CacheInvalidator(redisCache))
		// Add page caching middleware
		r.Use(customMiddleware.CachePage(redisCache, *cacheTTL))
	}

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Register routes
	h.Register(r)
	h.RegisterAPI(r)

	// Add editor route
	r.Get("/editor", func(w http.ResponseWriter, r *http.Request) {
		renderView(w, templates.EditorPage(*siteName))
	})

	// Set up the server
	server := &http.Server{
		Addr:         ":" + *port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Start the server
	log.Printf("Server started on port %s", *port)
	log.Printf("Blog available at %s", *baseURL)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

// Render templ views
func renderView(w http.ResponseWriter, component templ.Component) {
	err := component.Render(context.Background(), w)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Page templates
var (
	homePage      = templates.HomePage
	mePage        = templates.MePage
	portfolioPage = templates.PortfolioPage
	booksPage     = templates.BooksPage
	languagesPage = templates.LanguagesPage
	toolsPage     = templates.ToolsPage
	postPage      = templates.PostPage
	editorPage    = templates.EditorPage
) 