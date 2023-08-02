build:
	go build -o bin/main

run: build
	./bin/main

test:
	go test -v ./...

tidy:
	go fmt
	go mod tidy

.PHONY: build, run, test, tidy