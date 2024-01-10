package types

import "time"

type Product struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsTracked bool      `json:"is_tracked"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProduct(name string, url string) *Product {
	return &Product{
		Name:      name,
		IsTracked: true,
		Url:       url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type ProductVariant struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Url   string `json:"url"`
}

func NewProductVariant(name string, price int, url string) *ProductVariant {
	return &ProductVariant{
		Name:  name,
		Price: price,
		Url:   url,
	}
}

type CreateProductRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
