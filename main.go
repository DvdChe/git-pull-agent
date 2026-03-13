package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DvdChe/git-pull-agent/gitclient"

	"github.com/DvdChe/git-pull-agent/config"
)

func main() {
	logger := log.New(os.Stdout, "[git-pull-agent] ", log.LstdFlags|log.Lshortfile)
	logger.Println("Git Pull Agent starting...")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Error loading configuration: %v", err)
	}

	logger.Printf("Configuration loaded: RepoURL=%s, SyncPath=%s, Interval=%s, AuthMethod=%s",
		cfg.RepoURL, cfg.SyncPath, cfg.Interval, cfg.AuthMethod)

	// Create a context that can be cancelled to gracefully shut down the agent
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		logger.Printf("Received signal %s, initiating graceful shutdown...", sig)
		cancel()
	}()

	// Initial clone or pull
	logger.Println("Performing initial Git operation...")
	if err := gitclient.CloneOrPull(cfg); err != nil {
		logger.Printf("Initial Git operation failed: %v", err)
		// Depending on desired behavior, could exit here or continue trying.
		// For now, we'll continue but log the error.
	}
	logger.Println("Initial Git operation completed.")

	// Start the periodic pull loop
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Printf("Performing periodic Git pull (interval: %s)...", cfg.Interval)
			if err := gitclient.CloneOrPull(cfg); err != nil {
				logger.Printf("Periodic Git pull failed: %v", err)
			} else {
				logger.Println("Periodic Git pull completed successfully.")
			}
		case <-ctx.Done():
			logger.Println("Agent stopping due to context cancellation.")
			return
		}
	}
}
