package storage

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mluksic/product-price-tracker/types"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
	"os"
)

type PostgresStorage struct {
	db *sql.DB
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewPostgresStorage() *PostgresStorage {
	dbUrl := fmt.Sprintf("%s", viper.Get("DB_URL"))
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//defer conn.Close(context.Background())
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	return &PostgresStorage{
		db: db,
	}
}

func (p PostgresStorage) GetProductPrices(pId int) ([]types.ProductPrice, error) {
	rows, err := p.db.Query("select pp.price, pp.fetched_at from product_price pp where product_id = $1", pId)
	if err != nil {
		return nil, err
	}

	prices := []types.ProductPrice{}
	for rows.Next() {
		var price types.ProductPrice
		err := rows.Scan(&price.Price, &price.FetchedAt)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)

		fmt.Println("prices found", price)
	}

	return prices, nil
}

func (p PostgresStorage) CreateProductPrice(productPrice *types.ProductPrice) error {
	_, err := p.db.Exec("insert into product_price(price, fetched_at, product_id) values ($1, $2, $3)", productPrice.Price, productPrice.FetchedAt, productPrice.ProductId)

	return err
}

func (p PostgresStorage) CreateProduct(product *types.Product) error {
	_, err := p.db.Exec("insert into product (name, is_tracked, url, created_at, updated_at) values ($1,$2,$3,$4,$5)", product.Name, product.IsTracked, product.Url, product.CreatedAt, product.UpdatedAt)

	return err
}

func (p PostgresStorage) DeleteProduct(id int) error {
	_, err := p.db.Exec("delete from product where id = $1", id)

	return err
}

func (p PostgresStorage) GetProducts() ([]types.Product, error) {
	rows, err := p.db.Query("select id, name, is_tracked, url, created_at, updated_at from product order by id desc")

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
	row, err := p.db.Query("select id, name, is_tracked, url, created_at, updated_at from product where id = $1", id)
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
	_, err := p.db.Exec("update product set is_tracked = not is_tracked where id = $1", id)

	return err
}

func (p PostgresStorage) DeleteProductPrice() error {
	//TODO implement me
	panic("implement me")
}
