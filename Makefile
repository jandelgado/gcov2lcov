# makefile for gcov2lcov
.PHONY: build test inttest

build:
	go build -o bin/gcov2lcov .

test:
	go test ./... -coverprofile coverage.out
	go tool cover -func coverage.out

inttest: build
	./bin/gcov2lcov -infile testdata/coverage.out -outfile coverage.lcov
	diff -y testdata/coverage_expected.lcov coverage.lcov

