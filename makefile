build:
	go build -o build/noimportsasvar ./cmd/noimportsasvar

test:
	go test ./...

lint:
	golangci-lint run ./...