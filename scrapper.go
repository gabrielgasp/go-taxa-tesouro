package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/gabrielgasp/go-taxa-tesouro/model"
	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
)

type Scrapper interface {
	Run(context.Context)
}

type scrapper struct {
	logger  *slog.Logger
	rwMutex *sync.RWMutex
	wg      *sync.WaitGroup
}

func NewScrapper(logger *slog.Logger, rxMutex *sync.RWMutex, wg *sync.WaitGroup) Scrapper {
	return scrapper{
		logger:  logger,
		rwMutex: rxMutex,
		wg:      wg,
	}
}

func (s scrapper) Run(ctx context.Context) {
	defer s.wg.Done()

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		s.logger.Error("Error loading location")
		return
	}

	s.scrape()
	s.logger.Info("Initial scrapping finished")

	ticker := time.NewTicker(viper.GetDuration("INTERVAL_MINUTES") * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if s.shouldScrape(loc) {
				s.logger.Debug("Scrapping started")
				s.scrape()
				s.logger.Debug("Scrapping finished")
			}
		case <-ctx.Done():
			s.logger.Info("Scrapper stopped")
			return
		}
	}
}

func (s scrapper) shouldScrape(loc *time.Location) bool {
	now := time.Now().In(loc)

	return now.Weekday() >= time.Weekday(viper.GetInt("START_DAY")) &&
		now.Weekday() <= time.Weekday(viper.GetInt("END_DAY")) &&
		now.Hour() >= viper.GetInt("START_HOUR") &&
		now.Hour() < viper.GetInt("END_HOUR")
}

func (s scrapper) scrape() {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	fakeChrome := req.ImpersonateChrome()
	res := fakeChrome.Get(viper.GetString("URL_TESOURO")).Do()

	if res.StatusCode != 200 {
		s.logger.Error("Failed to fetch data", "status code", res.StatusCode)
		return
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", "error", err.Error())
		return
	}

	var tesouroResponse model.TesouroResponse
	if err := json.Unmarshal(data, &tesouroResponse); err != nil {
		s.logger.Error("Failed to unmarshal response", "error", err.Error())
		return
	}

	scrapperCache.Save(tesouroResponse.Data)
}
