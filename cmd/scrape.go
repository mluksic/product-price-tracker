package cmd

import (
	"fmt"
	"github.com/mluksic/product-price-tracker/scraper"
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/spf13/cobra"
	"os"
	"slices"
	"strings"
)

var Site string
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape product prices",
	Long:  "Scrape products prices from specific source/store",
	Run: func(cmd *cobra.Command, args []string) {
		store := storage.NewPostgresStorage()
		scrapers := []string{
			"nepremicnine",
			"amazon",
		}
		if !slices.Contains(scrapers, Site) {
			fmt.Printf("You have to one of supported scrapers: %s", strings.Join(scrapers, ","))
			os.Exit(0)
		}

		if Site == "nepremicnine" {
			scraper := scraper.NewNepremicnineScraper(store)
			scraper.ScrapeAndSave()
		} else {
			scraper := scraper.NewAmazonScraper(store)
			scraper.ScrapeAndSave()
		}
	},
}

func init() {
	scrapeCmd.Flags().StringVarP(&Site, "site", "s", "", "Select which scraper to use - 'nepremicnine' or 'amazon'")
	scrapeCmd.MarkFlagRequired("site")
	rootCmd.AddCommand(scrapeCmd)
}
