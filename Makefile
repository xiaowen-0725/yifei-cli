.PHONY: build test sync-data release tidy

build:
	CGO_ENABLED=0 go build -o bin/yifei .

test:
	go test ./...

tidy:
	go mod tidy

# Copies metadata from the docs repo into internal/assets (see Task 5/6).
sync-data:
	cp ../yifei-erp-docs/schema.json internal/assets/schema.json

release:
	goreleaser release --clean
