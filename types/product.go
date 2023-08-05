package types

type Product struct {
	Name      string `json:"name"`
	IsTracked bool   `json:"is_tracked"`
}

func NewProduct(name string, isTracked bool) *Product {
	return &Product{
		Name:      name,
		IsTracked: isTracked,
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
}
