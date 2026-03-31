all: build test

build:
	go build ./...

test: build
	go test ./...

clean:
	rm -rf build/
