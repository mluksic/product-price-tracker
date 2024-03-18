package api

import (
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/mluksic/product-price-tracker/handlers"
	"github.com/mluksic/product-price-tracker/scraper"
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/mluksic/product-price-tracker/types"
	"github.com/mluksic/product-price-tracker/views"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Config Config
}

func NewServer(config Config) *Server {
	return &Server{
		Config: config,
	}
}

type Config struct {
	Id         int
	Name       string
	ListenAddr string
	Storage    storage.Storer
}

func NewConfig() Config {
	return Config{
		Id:         12,
		Name:       "new server",
		ListenAddr: ":3030",
		Storage:    storage.NewPostgresStorage(),
	}
}

func (c Config) WithId(id int) Config {
	c.Id = id
	return c
}

func (c Config) WithName(name string) Config {
	c.Name = name
	return c
}

func (c Config) WithListenAddr(addr string) Config {
	c.ListenAddr = addr
	return c
}

func (c Config) WithStorage(s storage.Storer) Config {
	c.Storage = s
	return c
}
func (s *Server) Start() error {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	r.Get("/", s.handleIndexPage)

	// handlers
	productHandler := handlers.NewProductHandler(s.Config.Storage)

	// routes for "products" resource
	r.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.HandleIndex)
		r.Post("/", s.handleCreateProduct)

		// routes for "products/{Id}"
		r.Route("/{Id}", func(r chi.Router) {
			r.Get("/", s.handleGetProductPrices)
			r.Post("/scrape", s.handleScrapeProductPrices)
			r.Put("/track", s.handleToggleProductTracking)
			r.Delete("/", s.handleProductDeletion)
		})
	})

	return http.ListenAndServe(s.Config.ListenAddr, r)
}

func (s *Server) handleGetProductPrices(w http.ResponseWriter, r *http.Request) {
	productId, err := getId(r)
	if err != nil {
		views.ItemCreatedAlert(false, fmt.Sprintf("Unable to fetch product prices: %s", err.Error())).Render(r.Context(), w)
		return
	}
	prices, err := s.Config.Storage.GetProductPrices(productId)
	if err != nil {
		views.ItemCreatedAlert(false, fmt.Sprintf("Unable to fetch product prices: %s", err.Error())).Render(r.Context(), w)
		return
	}

	views.ProductPricesTable(prices).Render(r.Context(), w)
}

func (s *Server) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		views.ItemCreatedAlert(false, fmt.Sprintf("There was an issue parsing the form - %s", err.Error())).Render(r.Context(), w)
		return
	}

	p := types.NewProduct(r.PostFormValue("name"), r.PostFormValue("url"))
	err = s.Config.Storage.CreateProduct(p)
	if err != nil {
		views.ItemCreatedAlert(false, fmt.Sprintf("Unable to create new product price in the DB - %s", err.Error())).Render(r.Context(), w)
		return
	}

	w.Header().Set("Hx-Trigger", "product-added")
	views.ItemCreatedAlert(true, "You are successfully tracking new product price").Render(r.Context(), w)
}

func (s *Server) handleProductDeletion(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ApiError{
			Error:  "There was an error fetching query param",
			Detail: err.Error(),
		})
		return
	}

	err = s.Config.Storage.DeleteProduct(id)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ApiError{
			Error:  "Unable to delete product",
			Detail: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (s *Server) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	products, err := s.Config.Storage.GetProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templ.Handler(views.Show(products)).ServeHTTP(w, r)
}

func (s *Server) handleScrapeProductPrices(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ApiError{
			Error:  "There was an error fetching query param",
			Detail: err.Error(),
		})
		return
	}
	product, err := s.Config.Storage.GetProduct(id)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ApiError{
			Error:  "There was an error retrieving the product from the DB",
			Detail: err.Error(),
		})
		return
	}

	var targetScraper scraper.Scraper
	if strings.Contains(product.Url, "amazon") {
		targetScraper = scraper.NewAmazonScraper(s.Config.Storage)
	} else if strings.Contains(product.Url, "nepremicnine") {
		targetScraper = scraper.NewNepremicnineScraper(s.Config.Storage)
	} else if strings.Contains(product.Url, "mimovrste") {
		targetScraper = scraper.NewMimovrsteScraper(s.Config.Storage)
	}

	productVariants, err := targetScraper.Scrape([]string{product.Url})
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, ApiError{
			Error:  "There was an error scraping the product",
			Detail: err.Error(),
		})
		return
	}

	// save scraped products into DB
	for _, productVariant := range productVariants {
		productPrice := types.NewProductPrice(product.ID, productVariant.Price, time.Now())
		err := s.Config.Storage.CreateProductPrice(productPrice)
		if err != nil {
			WriteJson(w, http.StatusInternalServerError, ApiError{
				Error:  "There was an saving scraped prices for product into the DB",
				Detail: err.Error(),
			})
			return
		}
	}

	WriteJson(w, http.StatusOK, map[string]string{"message": "successfully scraped product prices"})
}

func (s *Server) handleToggleProductTracking(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ApiError{
			Error:  "There was an error fetching query param",
			Detail: err.Error(),
		})
		return
	}

	err = s.Config.Storage.ToggleProductTracking(id)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ApiError{
			Error:  "There was an error toggling product tracking",
			Detail: err.Error(),
		})
		return
	}
	templ.Handler(views.ItemCreatedAlert(true, "Your action has been saved")).ServeHTTP(w, r)
}

func WriteJson(w http.ResponseWriter, status int, msg any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		http.Error(w, "Unable to write JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

type ApiError struct {
	Error  string `json:"error"`
	Detail string `json:"detail"`
}

func getId(r *http.Request) (int, error) {
	param := chi.URLParam(r, "Id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return id, fmt.Errorf("invalid Id param given %s", param)
	}

	return id, nil
}
