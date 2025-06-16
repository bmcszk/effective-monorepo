package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/bmcszk/effective-monorepo/services/cleanup-sidecar/internal/cleanup"
	"github.com/bmcszk/effective-monorepo/services/cleanup-sidecar/internal/config"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("starting cleanup sidecar", "config", cfg)

	cleaner, err := cleanup.New(cfg)
	if err != nil {
		slog.Error("failed to create cleaner", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := cleaner.Close(); err != nil {
			slog.Error("failed to close cleaner", "error", err)
		}
	}()

	if err := cleaner.Run(ctx); err != nil {
		slog.Error("cleanup failed", "error", err)
		os.Exit(1)
	}

	slog.Info("cleanup completed successfully")
	fmt.Println("Cleanup sidecar completed successfully")
}
