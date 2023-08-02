package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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
	defer conn.Close(context.Background())

	return &PostgresStorage{
		db: conn,
	}
}

func (db *PostgresStorage) Get() {

}
func (db *PostgresStorage) Insert() {

}
func (db *PostgresStorage) Update() {

}
func (db *PostgresStorage) Delete() {

}
