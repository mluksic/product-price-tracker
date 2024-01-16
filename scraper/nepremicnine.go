package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/mluksic/product-price-tracker/types"
	"github.com/mluksic/product-price-tracker/util"
	"log"
	"strings"
	"time"
)

type NepremicnineScraper struct {
	c       *colly.Collector
	storage storage.Storer
}

func NewNepremicnineScraper(storage storage.Storer) *NepremicnineScraper {
	return &NepremicnineScraper{
		c:       colly.NewCollector(),
		storage: storage,
	}
}

func (scraper NepremicnineScraper) Scrape(urls []string) ([]types.ProductVariant, error) {
	var products []types.ProductVariant

	scraper.c.OnHTML(".cena span", func(e *colly.HTMLElement) {
		priceStr := strings.Replace(e.Text, "cca", "", -1)
		priceStr = strings.Replace(priceStr, "€", "", -1)

		product := types.ProductVariant{Price: util.PriceToCents(priceStr)}
		products = append(products, product)
	})

	// Error handling
	scraper.c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", string(r.Body), "\nError:", err)
	})

	// Print for debugging purposes
	scraper.c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: " + r.URL.String())
	})

	for _, url := range urls {
		err := scraper.c.Visit(url)
		if err != nil {
			return nil, err
		}
	}

	return products, nil
}

func (n NepremicnineScraper) ScrapeAndSave() {
	products, err := n.storage.GetProducts()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, product := range products {
		if !product.IsTracked {
			continue
		}

		if !strings.Contains(product.Url, "nepremicnine") {
			fmt.Println("skipping...")
			continue
		}

		productVariants, err := n.Scrape([]string{product.Url})
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, variant := range productVariants {
			price := types.ProductPrice{
				ProductId: product.ID,
				Price:     variant.Price,
				FetchedAt: time.Now(),
			}

			err := n.storage.CreateProductPrice(&price)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
