build:
	go build -o bin/gcov2lcov .

test:
	go test ./... -coverprofile coverage.out
	go tool cover -func coverage.out
