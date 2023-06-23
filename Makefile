default: build
.PHONY: build
build:
	go build -o dist/funcguard cmd/funcguard/main.go

.PHONY: test
test:
	go test -race ./...

.PHONY: check
check:
	golangci-lint run
