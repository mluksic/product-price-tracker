package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/mluksic/product-price-tracker/types"
	"github.com/mluksic/product-price-tracker/util"
	"html/template"
	"net/http"
	"strconv"
)

type Server struct {
	listenAddr string
	storage    storage.Storer
}

func NewServer(listenAddr string, store storage.Storer) *Server {
	return &Server{
		listenAddr: listenAddr,
		storage:    store,
	}
}

func (s *Server) Start() error {
	r := chi.NewRouter()

	r.Get("/", s.handleIndexPage)

	r.Get("/products", s.handleGetProducts)
	r.Post("/products", s.handleCreateProduct)
	r.Get("/products/{id}", s.handleGetProductPrices)

	http.HandleFunc("/", s.handleIndexPage)

	return http.ListenAndServe(s.listenAddr, r)
}

func (s *Server) handleGetProductPrices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productId, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	prices, err := s.storage.GetProductPrices(productId)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	err = json.NewEncoder(w).Encode(prices)
	if err != nil {
		http.Error(w, "Unable to encode response to JSON", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req types.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := types.NewProduct(req.Name)
	err = s.storage.CreateProduct(p)
	if err != nil {
		http.Error(w, "Unable to create product "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products, err := s.storage.GetProducts()
	if err != nil {
		http.Error(w, "Unable to get products: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(&products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Server) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	var tmplFile = "index.html"
	// add template functions
	funcMap := template.FuncMap{
		"IntToFloat": util.IntToFloat,
	}

	tmpl, err := template.New(tmplFile).Funcs(funcMap).ParseFiles("templates/" + tmplFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pPrices, _ := s.storage.GetProductPrices(1)
	products, _ := s.storage.GetProducts()

	tmplData := map[string]any{
		"product_prices": pPrices,
		"products":       products,
	}

	err = tmpl.Execute(w, tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getId(r *http.Request) (int, error) {
	param := chi.URLParam(r, "id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return id, fmt.Errorf("invalid id param given %s", param)
	}

	return id, nil
}
