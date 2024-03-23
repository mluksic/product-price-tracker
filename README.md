# Product Price Tracker

For tracking product prices from e-comm sites

Check [WIKI](https://github.com/mluksic/product-price-tracker/wiki) for more info

## Dependencies

- [Go](https://go.dev/doc/install)
- [go-chi](https://go-chi.io/#/pages/getting_started) router
- [Docker Compose](https://docs.docker.com/compose/install/) (optional) - for Postgres DB
- [Goose](https://github.com/pressly/goose) - for handling DB migrations

## Prerequisites

Download and install:

- [Go](https://go.dev/doc/install)
- [Goose](https://github.com/pressly/goose)


## Running the CLI app
App has two main commands. (run `make run` to get more info)
- serve - server HTTP serve on specific port (:3030 is default)
- scrape - fetch prices from URLs

### Create `.env` file in the root

Copy, rename, and adjust [.env.template](./.env.template) file so that the app connects to your Postgres DB

### Run the following command to start the CLI app
```bash
$ make run
```

### Run the following command to use live reload (requires [Air](https://github.com/cosmtrek/air) to be installed)
```bash
$ air
```

## Build

```bash
$ make build
```

## Tidy

```bash
$ make tidy
```

## Live reload

Project uses [Air](https://github.com/cosmtrek/air) for better & faster developer experience

## Using Goose to handle DB migrations

### Install [Goose](https://github.com/pressly/goose) on your machine

For Mac via Brew:
```bash
$ brew install goose
```

For convenience, you can set the Goose DB driver and DB string
```bash
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=postgresql://tracker:tracker@127.0.0.1:5432/tracker?sslmode=disable
```

### Check DB status

```bash
$ goose -dir storage/migrations status
```

### Create new migration file

```bash
$ goose -dir storage/migrations create create_product_price sql
```

### Run migration

```bash
$ goose -dir storage/migrations up
```

### Rollback migration

```bash
$ goose -dir storage/migrations down
```

## Deploy
The project uses Github Actions for CI/CD. Open [.deploy.yml](.github/workflows/deploy.yml) file to see the setup.

Push to `main` branch will automatically run CI/CD pipeline and deploy new changes to server.

## Test
(TBD)

## Contributors

👤 [Miha Luksic](https://www.mihaluksic.com)

