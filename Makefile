# makefile for gcov2lcov
.PHONY: phony

all: build-linux

phony:

build-linux: phony
	GOOS=linux GOARCH=amd64 go build -o bin/gcov2lcov-linux-amd64 .

build-windows: phony
	GOOS=windows GOARCH=amd64 go build -o bin/gcov2lcov-win-amd64 .

build-darwin: phony
	GOOS=darwin GOARCH=amd64 go build -o bin/gcov2lcov-darwin-amd64 .

test: phony
	go test ./... -coverprofile coverage.out
	go tool cover -func coverage.out

inttest: phony
	./bin/gcov2lcov-linux-amd64 -infile testdata/coverage.out -outfile coverage.lcov
	diff -y testdata/coverage_expected.lcov coverage.lcov

clean: phony
	rm -f bin/*

