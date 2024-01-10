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

type AmazonScraper struct {
	c       *colly.Collector
	storage storage.Storer
}

func NewAmazonScraper(store storage.Storer) *AmazonScraper {
	return &AmazonScraper{
		c:       colly.NewCollector(),
		storage: store,
	}
}

func (scraper AmazonScraper) Scrape(urls []string) ([]types.ProductVariant, error) {
	var (
		products []types.ProductVariant
	)
	scraper.c.OnHTML("span.a-price.aok-align-center.reinventPricePriceToPayMargin.priceToPay", func(e *colly.HTMLElement) {
		priceWhole := e.ChildText(".a-price-whole")
		priceFraction := e.ChildText(".a-price-fraction")
		priceStr := priceWhole + priceFraction
		priceStr = strings.Replace(priceStr, ".", "", -1)
		price, _ := strconv.Atoi(priceStr)
		product := types.ProductVariant{
			Price: price,
		}
		products = append(products, product)
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

//func (scraper AmazonScraper) getSearchUrls(productNames []string) []string {
//	var urls []string
//
//	for _, pName := range productNames {
//		queryParam := strings.Join(strings.Split(pName, " "), "+")
//		url := "https://www.amazon.de/s?k=" + queryParam
//		urls = append(urls, url)
//	}
//
//	return urls
//}

func (scraper AmazonScraper) ScrapeAndSave() {
	products, err := scraper.storage.GetProducts()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, product := range products {
		// skip disabled product
		if !product.IsTracked {
			continue
		}
		productVariants, err := scraper.Scrape([]string{product.Url})
		if err != nil {
			log.Fatal(err.Error())
		}

		// save scraped products into DB
		for _, productVariant := range productVariants {
			productPrice := types.NewProductPrice(product.ID, productVariant.Price, time.Now())
			err := scraper.storage.CreateProductPrice(productPrice)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}

	fmt.Println("Scraped and saved product prices to the DB...")
}
