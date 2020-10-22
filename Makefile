.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: migrate_up
migrate_up:
	go run -v ./cmd/apiserver migrate up

.PHONY: migrate_down
migrate_down:
	go run -v ./cmd/apiserver migrate down

.PHONY: serve
serve:
	go run -v ./cmd/apiserver serve

.DEFAULT_GOAL := build