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

	return http.ListenAndServe(s.listenAddr, r)
}

func (s *Server) handleGetProductPrices(w http.ResponseWriter, r *http.Request) {
	productId, err := getId(r)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		return
	}
	prices, err := s.storage.GetProductPrices(productId)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Unable to fetch product prices"})
		return
	}

	var tmplFile = "product_prices.html"
	// add template functions
	funcMap := template.FuncMap{
		"IntToFloat": util.IntToFloat,
	}

	tmpl, err := template.New(tmplFile).Funcs(funcMap).ParseFiles("templates/" + tmplFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplData := map[string]any{
		"prices": prices,
	}

	err = tmpl.Execute(w, tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	//WriteJson(w, http.StatusOK, prices)
}

func (s *Server) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req types.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		return
	}

	p := types.NewProduct(req.Name)
	err = s.storage.CreateProduct(p)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Unable to create product: " + err.Error()})
		return
	}

	WriteJson(w, http.StatusCreated, p)
}

func (s *Server) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := s.storage.GetProducts()
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ApiError{Error: "Unable to get products: " + err.Error()})
		return
	}

	WriteJson(w, http.StatusOK, products)
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

	products, _ := s.storage.GetProducts()
	tmplData := map[string]any{
		"products": products,
	}

	err = tmpl.Execute(w, tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	Error string `json:"error"`
}

func getId(r *http.Request) (int, error) {
	param := chi.URLParam(r, "id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return id, fmt.Errorf("invalid id param given %s", param)
	}

	return id, nil
}
