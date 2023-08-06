package api

import (
	"encoding/json"
	"github.com/mluksic/product-price-tracker/storage"
	"github.com/mluksic/product-price-tracker/types"
	"github.com/mluksic/product-price-tracker/util"
	"html/template"
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

	http.HandleFunc("/", s.handleIndexPage)

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

	products, _ := s.storage.GetProductPrices(1)
	err = tmpl.Execute(w, products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
