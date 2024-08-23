package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gabrielgasp/go-taxa-tesouro/model"
	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
)

var scrapperCache model.ScrapperCache

type Scrapper interface {
	Run(context.Context)
}

type scrapper struct {
	rwMutex *sync.RWMutex
	wg      *sync.WaitGroup
}

func NewScrapper(rxMutex *sync.RWMutex, wg *sync.WaitGroup) Scrapper {
	return scrapper{
		rwMutex: rxMutex,
		wg:      wg,
	}
}

func (s scrapper) Run(ctx context.Context) {
	defer s.wg.Done()

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Println("Error loading location")
		return
	}

	s.scrap()
	log.Println("Initial scrapping finished")

	ticker := time.NewTicker(viper.GetDuration("INTERVAL_MINUTES") * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if s.shouldScrap(loc) {
				log.Println("Scrapping started")
				s.scrap()
				log.Println("Scrapping finished")
			}
		case <-ctx.Done():
			log.Println("Scrapper shutdown")
			return
		}
	}
}

func (s scrapper) shouldScrap(loc *time.Location) bool {
	now := time.Now().In(loc)

	return now.Weekday() >= time.Weekday(viper.GetInt("START_DAY")) &&
		now.Weekday() <= time.Weekday(viper.GetInt("END_DAY")) &&
		now.Hour() >= viper.GetInt("START_HOUR") &&
		now.Hour() <= viper.GetInt("END_HOUR")
}

func (s scrapper) scrap() {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	fakeChrome := req.ImpersonateChrome()
	res := fakeChrome.Get(viper.GetString("URL_TESOURO")).Do()

	if res.StatusCode != 200 {
		log.Println("Failed to fetch data; Status code:", res.StatusCode)
		return
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return
	}

	var tesouroResponse model.TesouroResponse
	if err := json.Unmarshal(data, &tesouroResponse); err != nil {
		log.Println("Failed to unmarshal response:", err)
		return
	}

	scrapperCache.BondsList = make([]model.ScrapperBond, len(tesouroResponse.Response.TreasureBondsList))
	scrapperCache.BondsMap = make(map[string]model.ScrapperBond)

	for i, tb := range tesouroResponse.Response.TreasureBondsList {
		scrapperBond := model.ScrapperBond{
			Name:                    tb.TreasureBond.Nm,
			AnnualInvestmentRate:    tb.TreasureBond.AnulInvstmtRate,
			UnitaryInvestmentValue:  tb.TreasureBond.UntrInvstmtVal,
			MinimumInvestmentAmount: tb.TreasureBond.MinInvstmtAmt,
			AnnualRedemptionRate:    tb.TreasureBond.AnulRedRate,
			UnitaryRedemptionValue:  tb.TreasureBond.UntrRedVal,
			MinimumRedemptionValue:  tb.TreasureBond.MinRedVal,
		}

		scrapperCache.BondsList[i] = scrapperBond
		scrapperCache.BondsMap[strings.ToLower(scrapperBond.Name)] = scrapperBond
	}
	scrapperCache.UpdatedAt = tesouroResponse.Response.BusinessStatus.DtTm
}
