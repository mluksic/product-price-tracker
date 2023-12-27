package scraper

import (
	"github.com/mluksic/product-price-tracker/types"
	"log/slog"
	"time"
)

type Scraper interface {
	Scrape(productNames []string) ([]types.ProductVariant, error)
	ScrapeAndSave()
}

type PriceManager struct {
	scraper Scraper
}

func NewPriceManager(scraper Scraper) *PriceManager {
	return &PriceManager{
		scraper: scraper,
	}
}

func (priceManager PriceManager) RunPriceManagerPeriodically() {
	slog.Info("Started running scraper")
	ticker := time.NewTicker(24 * time.Hour)
	//defer ticker.Stop()

	// Creating channel using make
	tickerChan := make(chan bool)

	go func() {
		for {
			select {
			case <-tickerChan:
				return
			case <-ticker.C:
				priceManager.scraper.ScrapeAndSave()
			}
		}
	}()
}
