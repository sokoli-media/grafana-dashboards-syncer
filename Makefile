tests:
	go test ./...

format:
	golangci-lint run -v
