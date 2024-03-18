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

tailwind:
	@npx tailwindcss -i views/css/styles.css -o public/styles.css --watch

templ:
	@templ generate --proxy=http://localhost:3030 --watch

.PHONY: build, run, test, tidy