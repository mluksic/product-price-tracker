package types

type Product struct {
	Name  string
	Price int
	Url   string
}

func NewProduct(name string, price int, url string) *Product {
	return &Product{
		Name:  name,
		Price: price,
		Url:   url,
	}
}

type ProductPrice struct {
	Name      string
	ProductId int
	Price     int
	FetchedAt string
}

func NewProductPrice(name string, productId int, price int, fetchedAt string) *ProductPrice {
	return &ProductPrice{
		Name:      name,
		ProductId: productId,
		Price:     price,
		FetchedAt: fetchedAt,
	}
}
