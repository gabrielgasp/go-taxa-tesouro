package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/gabrielgasp/go-taxa-tesouro/model"
	"github.com/spf13/viper"
)

var scraperCache model.ScraperCache

func main() {
	if err := bootstrapConfig(); err != nil {
		fmt.Println("Failed to load config file, should 'ENV' be set to production?")
		os.Exit(1)
	}

	logger := bootstrapLogger()

	wg := &sync.WaitGroup{}

	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelCtx()

	scraper := NewScraper(logger, wg)
	api := NewApi(logger, wg)

	if viper.GetBool("SCRAPER_ENABLED") {
		wg.Add(1)
		go scraper.Run(ctx)
	}

	wg.Add(1)
	go api.Run(ctx)

	wg.Wait()

	logger.Info("App shutdown complete")
}

func bootstrapConfig() error {
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if strings.ToLower(viper.GetString("ENV")) != "production" {
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	os.Setenv("TZ", viper.GetString("TZ"))
	viper.SetDefault("RATE_LIMIT_PER_MINUTE", 10)

	return nil
}

func bootstrapLogger() *slog.Logger {
	var level slog.Level
	switch strings.ToLower(viper.GetString("LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	loggerOpts := slog.HandlerOptions{Level: level}
	return slog.New(slog.NewJSONHandler(os.Stderr, &loggerOpts))
}
