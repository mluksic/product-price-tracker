package storage

import "github.com/mluksic/product-price-tracker/types"

type Storer interface {
	GetProductPrices(pId int) ([]types.ProductPrice, error)
	CreateProduct(productPrice *types.Product) error
	CreateProductPrice(productPrice *types.ProductPrice) error
	DeleteProductPrice() error
}
