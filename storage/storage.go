package storage

import "github.com/mluksic/product-price-tracker/types"

type Storer interface {
	GetProductPrices(pId int) ([]types.ProductPrice, error)
	CreateProductPrice() error
	DeleteProductPrice() error
}
