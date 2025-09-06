package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gabrielgasp/go-taxa-tesouro/model"
	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
)

type Scraper interface {
	Run(context.Context)
}

type scraper struct {
	logger  *slog.Logger
	rwMutex *sync.RWMutex
	wg      *sync.WaitGroup
}

func NewScraper(logger *slog.Logger, rxMutex *sync.RWMutex, wg *sync.WaitGroup) Scraper {
	return scraper{
		logger:  logger,
		rwMutex: rxMutex,
		wg:      wg,
	}
}

func (s scraper) Run(ctx context.Context) {
	defer s.wg.Done()

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		s.logger.Error("Error loading location")
		return
	}

	s.scrape()
	s.logger.Info("Initial scraping finished")

	ticker := time.NewTicker(viper.GetDuration("INTERVAL_MINUTES") * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if s.shouldScrape(loc) {
				s.logger.Debug("Scraping started")
				s.scrape()
				s.logger.Debug("Scraping finished")
			}
		case <-ctx.Done():
			s.logger.Info("Scraper stopped")
			return
		}
	}
}

func (s scraper) shouldScrape(loc *time.Location) bool {
	now := time.Now().In(loc)

	return now.Weekday() >= time.Weekday(viper.GetInt("START_DAY")) &&
		now.Weekday() <= time.Weekday(viper.GetInt("END_DAY")) &&
		now.Hour() >= viper.GetInt("START_HOUR") &&
		now.Hour() < viper.GetInt("END_HOUR")
}

func (s scraper) scrape() {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	investData, err := s.fetchData(viper.GetString("INVEST_CSV_URL"))
	if err != nil {
		s.logger.Error("Failed to fetch invest data", "error", err.Error())
		return
	}

	var investDataParsed []model.Invest
	if err := model.ParseCSV(investData, &investDataParsed); err != nil {
		s.logger.Error("Failed to parse invest data", "error", err.Error())
		return
	}

	redeemData, err := s.fetchData(viper.GetString("REDEEM_CSV_URL"))
	if err != nil {
		s.logger.Error("Failed to fetch redeem data", "error", err.Error())
		return
	}

	var redeemDataParsed []model.Redeem
	if err := model.ParseCSV(redeemData, &redeemDataParsed); err != nil {
		s.logger.Error("Failed to parse redeem data", "error", err.Error())
		return
	}

	scraperCache.Save(investDataParsed, redeemDataParsed)
}

func (s scraper) fetchData(url string) ([]byte, error) {
	res, err := req.ImpersonateChrome().R().Get(url)
	if err != nil || res.Response == nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	return data, nil
}
