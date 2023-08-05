package types

type Product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Url   string `json:"url"`
}

func NewProduct(name string, price int, url string) *Product {
	return &Product{
		Name:  name,
		Price: price,
		Url:   url,
	}
}

type CreateProductRequest struct {
	Name string `json:"name"`
}
