package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Server represents an HTTP server
type Server struct {
	server *http.Server
}

// NewServer creates a new HTTP server with the given handler and port
func NewServer(handler http.Handler, port string) *Server {
	return &Server{
		server: &http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

// Start starts the HTTP server with graceful shutdown
func (s *Server) Start() error {
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
		log.Println("Shutting down server...")
		err := s.server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
		serverStopCtx()
	}()

	// Start the server
	log.Printf("Server started on %s", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	return nil
}

// Shutdown shuts down the server immediately
func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.Background())
}
