package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

type PostgresStorage struct{}

func (db *PostgresStorage) Connect() {
	dbUrl := "postgres://tracker:tracker@localhost:5432/tracker"
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
