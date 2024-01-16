package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/mluksic/product-price-tracker/types"
	"log"
	"strconv"
	"strings"
	"time"
)

type MimovrsteScraper struct {
	c       *colly.Collector
	storage storage.Storer
}

func NewMimovrsteScraper(storage storage.Storer) *MimovrsteScraper {
	return &MimovrsteScraper{
		c:       colly.NewCollector(),
		storage: storage,
	}
}

func (scraper MimovrsteScraper) Scrape(urls []string) ([]types.ProductVariant, error) {
	var products []types.ProductVariant

	scraper.c.OnHTML(".price__wrap__box__final", func(e *colly.HTMLElement) {
		priceStr := strings.Replace(e.Text, "€", "", -1)
		priceStr = strings.ReplaceAll(priceStr, "\n", "")
		priceStr = strings.ReplaceAll(priceStr, " ", "")
		priceStr = strings.Replace(priceStr, "\u00A0", "", -1)
		price, _ := strconv.Atoi(priceStr)

		product := types.ProductVariant{Price: price * 100}
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

func (n MimovrsteScraper) ScrapeAndSave() {
	products, err := n.storage.GetProducts()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, product := range products {
		if !product.IsTracked {
			continue
		}

		if !strings.Contains(product.Url, "mimovrste") {
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
