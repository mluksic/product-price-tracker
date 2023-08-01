# Product Price Tracker

For tracking product prices from e-comm sites

Check [WIKI](https://github.com/mluksic/product-price-tracker/wiki) for more info

## Dependencies

- [Go](https://go.dev/doc/install)
- [Docker Compose](https://docs.docker.com/compose/install/) (optional) - for Postgres DB
- [Goose](https://github.com/pressly/goose) - for handling DB migrations

## Prerequisites

Download and install:

- [Go](https://go.dev/doc/install)
- [Goose](https://github.com/pressly/goose)
 

## Running the app

### Create `.env` file in the root

Copy, rename, and adjust [.env.template](./.env.template) file so that the app connects to your Postgres DB

### Run the following command to start the app
```bash
$ go run main.go
```

## Build

```bash
$ go build -o main
```

## Using Goose to handle DB migrations

### Install [Goose](https://github.com/pressly/goose) on your machine

For Mac via Brew:
```bash
$ brew install goose
```

### Check DB status

```bash
$ goose -dir storage/migrations postgres "postgresql://tracker:tracker@localhost:5432/tracker?sslmode=disable" status
```

### Create new migration file

```bash
$ goose -dir storage/migrations postgres "postgresql://tracker:tracker@localhost:5432/tracker?sslmode=disable" create create_product_price sql
```

### Run migration

```bash
$ goose -dir storage/migrations postgres "postgresql://tracker:tracker@localhost:5432/tracker?sslmode=disable" up
```

### Rollback migration

```bash
$ goose -dir storage/migrations postgres "postgresql://tracker:tracker@localhost:5432/tracker?sslmode=disable" down
```

## Deploy
(TBD)

## Test
(TBD)

## Contributors

👤 [Miha Luksic](https://www.mihaluksic.com)

