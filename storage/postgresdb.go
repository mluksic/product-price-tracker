package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"os"
)

type PostgresStorage struct{}

func NewPostgresStorage() *PostgresStorage {
	return &PostgresStorage{}
}

func (db *PostgresStorage) Connect() {
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
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}
func (db *PostgresStorage) Get() {

}
func (db *PostgresStorage) Insert() {

}
func (db *PostgresStorage) Update() {

}
func (db *PostgresStorage) Delete() {

}
