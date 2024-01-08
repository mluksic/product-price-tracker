package cmd

import (
	"github.com/mluksic/product-price-tracker/scraper"
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/spf13/cobra"
)

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape product prices",
	Long:  "Scrape products prices from specific source/store",
	Run: func(cmd *cobra.Command, args []string) {
		store := storage.NewPostgresStorage()
		scraper := scraper.NewAmazonScraper(store)

		scraper.ScrapeAndSave()
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
}
