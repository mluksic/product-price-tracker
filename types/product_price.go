package types

import "time"

type ProductPrice struct {
	ProductId int       `json:"product_id"`
	Price     int       `json:"price"`
	FetchedAt time.Time `json:"fetched_at"`
}

func NewProductPrice(productId int, price int, fetchedAt time.Time) *ProductPrice {
	return &ProductPrice{
		ProductId: productId,
		Price:     price,
		FetchedAt: fetchedAt,
	}
}
