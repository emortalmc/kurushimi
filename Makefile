lint:
	golangci-lint run

pre-commit:
	go mod tidy
	make lint
