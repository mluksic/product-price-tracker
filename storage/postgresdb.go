package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/mluksic/product-price-tracker/types"
	"github.com/spf13/viper"
	"os"
)

type PostgresStorage struct {
	db *pgx.Conn
}

func NewPostgresStorage() *PostgresStorage {
	dbUrl := fmt.Sprintf("%s", viper.Get("DB_URL"))
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//defer conn.Close(context.Background())

	return &PostgresStorage{
		db: conn,
	}
}

func (p PostgresStorage) GetProductPrices(pId int) ([]types.ProductPrice, error) {
	rows, err := p.db.Query(context.Background(), "select pp.name, pp.price, pp.fetched_at, pp.url from product_price pp where product_id = $1", pId)
	if err != nil {
		return nil, err
	}

	prices := []types.ProductPrice{}
	for rows.Next() {
		var price types.ProductPrice
		err := rows.Scan(&price.Name, &price.Price, &price.FetchedAt, &price.Url)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)

		fmt.Println("prices found", price)
	}

	return prices, nil
}

func (p PostgresStorage) CreateProductPrice(productPrice *types.ProductPrice) error {
	_, err := p.db.Exec(context.Background(), "insert into product_price(name, price, url, fetched_at, product_id) values ($1, $2, $3, $4, $5)", productPrice.Name, productPrice.Price, productPrice.Url, productPrice.FetchedAt, productPrice.ProductId)

	return err
}

func (p PostgresStorage) CreateProduct(product *types.Product) error {
	_, err := p.db.Exec(context.Background(), "insert into product (name, is_tracked, url, created_at, updated_at) values ($1,$2,$3,$4,$5)", product.Name, product.IsTracked, product.Url, product.CreatedAt, product.UpdatedAt)

	return err
}

func (p PostgresStorage) DeleteProduct(id int) error {
	_, err := p.db.Exec(context.Background(), "delete from product where id = $1", id)

	return err
}

func (p PostgresStorage) GetProducts() ([]types.Product, error) {
	rows, err := p.db.Query(context.Background(), "select id, name, is_tracked, url, created_at, updated_at from product order by id desc")

	if err != nil {
		return nil, err
	}

	var products []types.Product
	for rows.Next() {
		var product types.Product
		err := rows.Scan(&product.ID, &product.Name, &product.IsTracked, &product.Url, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (p PostgresStorage) GetProduct(id int) (types.Product, error) {
	row, err := p.db.Query(context.Background(), "select id, name, is_tracked, url, created_at, updated_at from product where id = $1", id)
	if err != nil {
		return types.Product{}, err
	}

	var product types.Product

	for row.Next() {
		err = row.Scan(&product.ID, &product.Name, &product.IsTracked, &product.Url, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return types.Product{}, err
		}
	}

	return product, nil
}

func (p PostgresStorage) ToggleProductTracking(id int) error {
	_, err := p.db.Exec(context.Background(), "update product set is_tracked = not is_tracked where id = $1", id)

	return err
}

func (p PostgresStorage) DeleteProductPrice() error {
	//TODO implement me
	panic("implement me")
}
