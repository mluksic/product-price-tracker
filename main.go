package main

import (
	"flag"
	"fmt"
	"github.com/mluksic/product-price-tracker/api"
	"github.com/mluksic/product-price-tracker/scraper"
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/mluksic/product-price-tracker/types"
	"github.com/spf13/viper"
	"log"
	"time"
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

	viper.SetConfigFile(".env")
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	listenAddr := flag.String("listenAddr", ":3000", "the server port")
	flag.Parse()

	store := storage.NewPostgresStorage()
	server := api.NewServer(*listenAddr, store)

	// save scraped products into DB
	for _, product := range products {
		productPrice := types.NewProductPrice(product.Name, 1, product.Price, time.Now())
		err := store.CreateProductPrice(productPrice)
		if err != nil {
			log.Fatal("Unable to create product price: " + err.Error())
		}
	}

	fmt.Println("Started server on port" + *listenAddr)
	err2 := server.Start()

	if err2 != nil {
		log.Fatal("Unable to start HTTP server: " + err2.Error())
	}
}
