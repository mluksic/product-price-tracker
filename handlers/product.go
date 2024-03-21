package handlers

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/mluksic/product-price-tracker/types"
	"github.com/mluksic/product-price-tracker/views"
	"net/http"
	"strconv"
)

func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	products, _ := h.s.GetProducts()

	templ.Handler(views.Show(products)).ServeHTTP(w, r)
}

func (h *Handler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		templ.Handler(views.ItemCreatedAlert(false, fmt.Sprintf("There was an issue parsing the form - %s", err.Error()))).ServeHTTP(w, r)
		return
	}

	p := types.NewProduct(r.PostFormValue("name"), r.PostFormValue("url"))
	err = h.s.CreateProduct(p)
	if err != nil {
		templ.Handler(views.ItemCreatedAlert(false, fmt.Sprintf("Unable to create new product price in the DB - %s", err.Error()))).ServeHTTP(w, r)
		return
	}
}

func (h *Handler) HandleProductDeletion(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "Id")

	id, err := strconv.Atoi(param)
	if err != nil {
		templ.Handler(views.ItemCreatedAlert(false, fmt.Sprintf("Unable to fetch request query param: %s", err.Error()))).ServeHTTP(w, r)
		return
	}

	err = h.s.DeleteProduct(id)
	if err != nil {
		templ.Handler(views.ItemCreatedAlert(false, fmt.Sprintf("Unablet to delete product: %s", err.Error()))).ServeHTTP(w, r)
		return
	}

	templ.Handler(views.ItemCreatedAlert(true, "You've successfully deleted product")).ServeHTTP(w, r)
}
