package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/mluksic/product-price-tracker/handlers"
	"github.com/mluksic/product-price-tracker/storage"
	"net/http"
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

	// handlers
	handler := handlers.NewHandler(s.Config.Storage)

	// routes
	r.HandleFunc("/", handler.HandleIndexPage)

	// routes for "products" resource
	r.Route("/products", func(r chi.Router) {
		r.Get("/", handler.HandleIndex)
		r.Post("/", handler.HandleCreateProduct)

		// routes for "products/{Id}"
		r.Route("/{Id}", func(r chi.Router) {
			r.Get("/", handler.HandleGetProductPrices)
			r.Post("/scrape", handler.HandleScrapeProductPrices)
			r.Put("/track", handler.HandleToggleProductTracking)
			r.Delete("/", handler.HandleProductDeletion)
		})
	})

	return http.ListenAndServe(s.Config.ListenAddr, r)
}
