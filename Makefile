all: build test integration-test

test:
	go test

integration-test:
	./test.sh

build:
	go build
