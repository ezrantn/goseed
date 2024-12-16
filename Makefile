build:
	@go build -o bin/goseed examples/main.go

run: build
	@./bin/goseed

test: 
	@go test -v ./...

fmt:
	@go fmt ./...