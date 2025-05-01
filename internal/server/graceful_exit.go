package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

// gracefulExit sets up signal handling for graceful server shutdown
func gracefulExit(
	serverCtx context.Context,
	serverStopCtx context.CancelFunc,
	s *Server,
) {
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
}
