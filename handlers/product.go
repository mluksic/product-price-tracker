package handlers

import (
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/mluksic/product-price-tracker/views"
	"net/http"
)

type ProductHandler struct {
	S storage.Storer
}

func NewProductHandler(s storage.Storer) ProductHandler {
	return ProductHandler{
		S: s,
	}
}

func (h ProductHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	products, _ := h.S.GetProducts()

	views.Show(products).Render(r.Context(), w)
}
