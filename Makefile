.PHONY: all build test lint run clean

BINARY := ip_extractor

all: build

build:
	mkdir -p bin
	go build -ldflags='-s -w' -o ./bin/$(BINARY) cmd/$(BINARY)/main.go

test:
	TESTING_FILES_DIR="$(shell pwd)/test_data" go test -v ./... -race -timeout=1m

lint:
	golangci-lint run

run: build
	./bin/$(BINARY) --file-path=${FILE_PATH} --strategy=${STRATEGY}

clean:
	rm -rf ./bin
