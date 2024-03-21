package handlers

import (
	"github.com/a-h/templ"
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/mluksic/product-price-tracker/views"
	"net/http"
)

type Handler struct {
	s storage.Storer
}

func NewHandler(s storage.Storer) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	products, err := h.s.GetProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templ.Handler(views.Show(products)).ServeHTTP(w, r)
}
