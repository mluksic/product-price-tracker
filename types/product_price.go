package types

import "time"

type ProductPrice struct {
	Name      string    `json:"name"`
	ProductId int       `json:"product_id"`
	Price     int       `json:"price"`
	Url       string    `json:"url"`
	FetchedAt time.Time `json:"fetched_at"`
}

func NewProductPrice(name string, productId int, price int, url string, fetchedAt time.Time) *ProductPrice {
	return &ProductPrice{
		Name:      name,
		ProductId: productId,
		Price:     price,
		Url:       url,
		FetchedAt: fetchedAt,
	}
}
