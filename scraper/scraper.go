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

type Scraper struct {
	c       *colly.Collector
	storage storage.Storer
}

func NewScraper(storer storage.Storer) *Scraper {
	return &Scraper{
		c:       colly.NewCollector(),
		storage: storer,
	}
}

func RunScraperPeriodically() {
	ticker := time.NewTicker(1 * time.Hour)
	//defer ticker.Stop()

	// Creating channel using make
	tickerChan := make(chan bool)

	go func() {
		for {
			select {
			case <-tickerChan:
				return
			case <-ticker.C:
				ScrapeAndSave()
			}
		}
	}()
}

func ScrapeAndSave() {
	store := storage.NewPostgresStorage()
	scraper := NewScraper(store)

	products, err := scraper.storage.GetProducts()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, product := range products {
		// skip disabled product
		if !product.IsTracked {
			continue
		}
		productVariants, err := Scrape([]string{product.Name})
		if err != nil {
			log.Fatal(err.Error())
		}

		// save scraped products into DB
		for _, productVariant := range productVariants {
			productPrice := types.NewProductPrice(productVariant.Name, product.ID, productVariant.Price, productVariant.Url, time.Now())
			err := store.CreateProductPrice(productPrice)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}

	fmt.Println("Scraped and saved product prices to the DB...")
}

func Scrape(productNames []string) ([]types.ProductVariant, error) {
	var (
		products []types.ProductVariant
		urls     []string
	)

	store := storage.NewPostgresStorage()
	scraper := NewScraper(store)

	urls = getSearchUrls(productNames)

	scraper.c.OnHTML(".s-result-list .s-result-item", func(e *colly.HTMLElement) {
		e.ForEach("div.a-section.a-spacing-base", func(_ int, h *colly.HTMLElement) {
			product := extractProduct(h)
			products = append(products, product)
		})
	})

	// Error handling
	scraper.c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Print for debugging purposes
	scraper.c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: " + r.URL.String())
	})

	// Start scraping URLs
	for _, url := range urls {
		err := scraper.c.Visit(url)
		if err != nil {
			return nil, err
		}
	}

	return products, nil
}

func extractProduct(h *colly.HTMLElement) types.ProductVariant {
	priceStr := strings.Join(strings.Split(h.ChildText("span.a-price-whole"), ","), "")
	price, _ := strconv.Atoi(priceStr)
	name := h.ChildText("span.a-text-normal")
	url := h.ChildAttr("a.a-link-normal", "href")

	return types.ProductVariant{
		Price: price,
		Url:   "https://www.amazon.de" + url,
		Name:  name,
	}
}

func getSearchUrls(productNames []string) []string {
	var urls []string

	for _, pName := range productNames {
		queryParam := strings.Join(strings.Split(pName, " "), "+")
		url := "https://www.amazon.de/s?k=" + queryParam
		urls = append(urls, url)
	}

	return urls
}
