package api

import (
	"encoding/json"
	"github.com/mluksic/product-price-tracker/storage"
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
