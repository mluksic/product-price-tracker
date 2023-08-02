package api

import (
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
	return http.ListenAndServe(s.listenAddr, nil)
}
