.PHONY: build run clean install fmt

BINARY_NAME=docktab
MAIN_PATH=./cmd/docktab

build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH) ps

clean:
	rm -f $(BINARY_NAME)

install:
	go install $(MAIN_PATH)

fmt:
	go fmt ./...