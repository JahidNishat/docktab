.PHONY: build run clean fmt

BINARY_NAME=docktab

build:
	go build -o $(BINARY_NAME) ./cmd/docktab

run:
	go run ./cmd/docktab ps

clean:
	rm -f $(BINARY_NAME)

fmt:
	go fmt ./...
