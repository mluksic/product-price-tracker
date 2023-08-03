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
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		viper.Get("DB_USERNAME"),
		viper.Get("DB_PASSWORD"),
		viper.Get("DB_HOST"),
		viper.Get("DB_PORT"),
		viper.Get("DB_DATABASE"),
	)
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
	rows, err := p.db.Query(context.Background(), "select pp.name, pp.price from product_price pp where product_id = $1", pId)
	if err != nil {
		return nil, err
	}

	prices := []types.ProductPrice{}
	for rows.Next() {
		var price types.ProductPrice
		err := rows.Scan(&price.Name, &price.Price)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)

		fmt.Println("prices found", price)
	}

	return prices, nil
}

func (p PostgresStorage) CreateProductPrice() error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresStorage) DeleteProductPrice() error {
	//TODO implement me
	panic("implement me")
}
