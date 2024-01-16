build:
	go build -o bin/main

run: build
	./bin/main -h

serve: build
	./bin/main serve -p :3030

scrape: build
	./bin/main scrape

test:
	go test -v ./...

tidy:
	go fmt
	go mod tidy

build-prod:
	env GOOS=linux GOARCH=amd64 go build -o bin/main

.PHONY: build, run, test, tidy