package main

import (
	"fmt"
	"github.com/mluksic/product-price-tracker/scraper"
	"log"
)

func main() {
	var productNames = []string{
		"apple tv",
		"google pixel 7",
	}

	products, err := scraper.Scrape(productNames)
	if err != nil {
		log.Fatal(err)
	}

	for _, product := range products {
		fmt.Println(product)
	}
}
