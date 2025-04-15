.DEFAULT_GOAL := build
.PHONY: clean


clean:
	rm -rf ./build
	go clean

build: clean
	go build -o ./build/app

run: build
	./build/app

test:
	go test -v ./...
