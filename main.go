package main

import (
	"flag"
	"fmt"
	"github.com/mluksic/product-price-tracker/api"
	"github.com/mluksic/product-price-tracker/scraper"
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	listenAddr := flag.String("listenAddr", ":3000", "the server port")
	flag.Parse()

	store := storage.NewPostgresStorage()
	server := api.NewServer(*listenAddr, store)

	scraper.RunScraperPeriodically()

	fmt.Println("Started server on port" + *listenAddr)
	err2 := server.Start()

	if err2 != nil {
		log.Fatal("Unable to start HTTP server: " + err2.Error())
	}
}
