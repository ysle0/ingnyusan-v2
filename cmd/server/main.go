package main

import (
	"flag"
	"log"

	"ingnyusan/v2/internal/cache"
	"ingnyusan/v2/internal/config"
	"ingnyusan/v2/internal/router"
	"ingnyusan/v2/internal/server"
	"ingnyusan/v2/internal/service"
)

func main() {
	// Define a flag for the config file path
	configPath := flag.String("config", "config.toml", "Path to TOML configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the cache
	redisCache := cache.SetupCache(cfg.Redis.URL)
	if redisCache != nil {
		defer redisCache.Close()
	}

	// Initialize the post service
	postService := service.NewPostService(cfg.Content.PostsDir)
	err = postService.LoadPosts()
	if err != nil {
		log.Fatalf("Failed to load posts: %v", err)
	}

	// Set up the router
	r := router.SetupRouter(cfg, postService, redisCache)

	// Create and start the server
	srv := server.NewServer(r, cfg.Server.Port)
	log.Printf("Blog available at %s", cfg.Server.BaseURL)
	err = srv.Start()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
