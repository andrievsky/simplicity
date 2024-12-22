.DEFAULT_GOAL := build
.PHONY: clean

build:
	rm -rf ./build
	go build -o ./build/app

run: build
	./build/app

test:
	go test -v ./...
