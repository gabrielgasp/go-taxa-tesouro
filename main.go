package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	if err := bootstrapConfig(); err != nil {
		log.Fatalf("Failed to load config file, should 'ENV' be set to production?")
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

	log.Println("App shutdown complete")
}
