package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
)

type Product struct {
	name  string
	price int
	url   string
}

type Scraper struct {
	c *colly.Collector
}

func newScraper() *Scraper {
	return &Scraper{
		c: colly.NewCollector(),
	}
}

func Scrape(productNames []string) ([]Product, error) {
	var (
		products []Product
		urls     []string
	)

	scraper := newScraper()

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

func extractProduct(h *colly.HTMLElement) Product {
	priceStr := strings.Join(strings.Split(h.ChildText("span.a-price-whole"), ","), "")
	price, _ := strconv.Atoi(priceStr)
	name := h.ChildText("span.a-text-normal")

	return Product{
		price: price,
		url:   h.Request.URL.String(),
		name:  name,
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