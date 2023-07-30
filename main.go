package main

import (
	"flag"
	"fmt"
	"github.com/mluksic/product-price-tracker/api"
	"log"
)

func main() {
	//var productNames = []string{
	//	"apple tv",
	//	"google pixel 7",
	//}
	//
	//products, err := scraper.Scrape(productNames)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, product := range products {
	//	fmt.Println(product)
	//}

	listenAddr := flag.String("listenAddr", ":3000", "the server port")
	flag.Parse()

	server := api.NewServer(*listenAddr)

	fmt.Println("Started server on port" + *listenAddr)
	err := server.Start()

	if err != nil {
		log.Fatal("Unable to start HTTP server: " + err.Error())
	}
}
