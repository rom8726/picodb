.PHONY: lint
lint:
	golangci-lint -v run

.PHONY: test
test:
	go test -v ./...
