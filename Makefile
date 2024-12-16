test: 
	@go test -v ./...

fmt:
	@go fmt ./...

pg-seed:
	@go run examples/postgres/main.go

mysql-seed:
	@go run examples/mysql/main.go