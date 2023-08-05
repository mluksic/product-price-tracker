package api

import (
	"encoding/json"
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/mluksic/product-price-tracker/types"
	"net/http"
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
	http.HandleFunc("/product_prices", s.handleGetProductPrices)
	http.HandleFunc("/products", s.handleCreateProduct)

	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *Server) handleGetProductPrices(w http.ResponseWriter, r *http.Request) {
	prices, err := s.storage.GetProductPrices(1)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	err = json.NewEncoder(w).Encode(prices)
	if err != nil {
		return
	}
}

func (s *Server) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	if "POST" != r.Method {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var p types.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(p.Name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
