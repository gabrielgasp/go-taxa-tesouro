package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	if err := bootstrapConfig(); err != nil {
		slog.Error("Failed to load config file, should 'ENV' be set to production?")
		os.Exit(1)
	}

	rwMutex := &sync.RWMutex{}
	wg := &sync.WaitGroup{}

	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelCtx()

	scrapper := NewScrapper(rwMutex, wg)
	api := NewApi(rwMutex, wg)

	wg.Add(2)
	go scrapper.Run(ctx)
	go api.Run(ctx)

	wg.Wait()

	slog.Info("App shutdown complete")
}
