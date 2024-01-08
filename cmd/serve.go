package cmd

import (
	"fmt"
	"github.com/mluksic/product-price-tracker/api"
	"github.com/mluksic/product-price-tracker/scraper"
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/spf13/cobra"
	"log"
)

var Port string
var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Start HTTP server",
	Long:    "Start HTTP server for an tracker app",
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		//listenAddr := flag.String("listenAddr", ":3000", "the server port")
		//flag.Parse()

		store := storage.NewPostgresStorage()
		server := api.NewServer(Port, store)

		amazonScraper := scraper.NewAmazonScraper(
			store,
		)
		priceManager := scraper.NewPriceManager(amazonScraper)

		priceManager.RunPriceManagerPeriodically()

		fmt.Println("Started server on port" + Port)
		err := server.Start()

		if err != nil {
			log.Fatal("Unable to start HTTP server: " + err.Error())
		}
	},
}

func init() {
	serveCmd.Flags().StringVarP(&Port, "port", "p", ":3030", "Define server port")
	rootCmd.AddCommand(serveCmd)
}
