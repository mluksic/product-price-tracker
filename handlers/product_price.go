package handlers

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/mluksic/product-price-tracker/scraper"
	"github.com/mluksic/product-price-tracker/types"
	"github.com/mluksic/product-price-tracker/views"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) HandleGetProductPrices(w http.ResponseWriter, r *http.Request) {
	productId, err := getId(r)
	if err != nil {
		templ.Handler(views.ItemCreatedAlert(false, fmt.Sprintf("Unable to fetch product prices: %s", err.Error()))).ServeHTTP(w, r)
		return
	}
	prices, err := h.s.GetProductPrices(productId)
	if err != nil {
		templ.Handler(views.ItemCreatedAlert(false, fmt.Sprintf("Unable to fetch product prices: %s", err.Error()))).ServeHTTP(w, r)
		return
	}

	templ.Handler(views.ProductPricesTable(prices)).ServeHTTP(w, r)
}

func (h *Handler) HandleScrapeProductPrices(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		templ.Handler(views.ItemCreatedAlert(false, fmt.Sprintf("Unable to fetch request query param: %s", err.Error()))).ServeHTTP(w, r)
		return
	}

	product, err := h.s.GetProduct(id)
	if err != nil {
		templ.Handler(views.ItemCreatedAlert(false, fmt.Sprintf("Unablet to fetch product: %s", err.Error()))).ServeHTTP(w, r)
		return
	}

	var targetScraper scraper.Scraper
	if strings.Contains(product.Url, "amazon") {
		targetScraper = scraper.NewAmazonScraper(h.s)
	} else if strings.Contains(product.Url, "nepremicnine") {
		targetScraper = scraper.NewNepremicnineScraper(h.s)
	} else if strings.Contains(product.Url, "mimovrste") {
		targetScraper = scraper.NewMimovrsteScraper(h.s)
	}

	productVariants, err := targetScraper.Scrape([]string{product.Url})
	if err != nil {
		templ.Handler(views.ItemCreatedAlert(false, fmt.Sprintf("Unable to scrape: %s", err.Error()))).ServeHTTP(w, r)
		return
	}

	// save scraped products into DB
	for _, productVariant := range productVariants {
		productPrice := types.NewProductPrice(product.ID, productVariant.Price, time.Now())
		err := h.s.CreateProductPrice(productPrice)

		if err != nil {
			templ.Handler(views.ItemCreatedAlert(false, fmt.Sprintf("Unable to create product price: %s", err.Error()))).ServeHTTP(w, r)
			return
		}
	}

	templ.Handler(views.ItemCreatedAlert(true, "You've successfully scraped product")).ServeHTTP(w, r)
}

func (h *Handler) HandleToggleProductTracking(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		templ.Handler(views.ItemCreatedAlert(false, fmt.Sprintf("Unable to fetch request query param: %s", err.Error()))).ServeHTTP(w, r)
		return
	}

	err = h.s.ToggleProductTracking(id)
	if err != nil {
		templ.Handler(views.ItemCreatedAlert(false, fmt.Sprintf("Unable to toggle tracking: %s", err.Error()))).ServeHTTP(w, r)
		return
	}

	templ.Handler(views.ItemCreatedAlert(true, "Your action has been saved")).ServeHTTP(w, r)
}

func getId(r *http.Request) (int, error) {
	param := chi.URLParam(r, "Id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return id, fmt.Errorf("invalid Id param given %s", param)
	}

	return id, nil
}
