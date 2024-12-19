test: 
	@go test -v ./...

fmt:
	@go fmt ./...

up:
	@sudo docker-compose up

pg-seed:
	@go run examples/postgres/main.go

mysql-seed:
	@go run examples/mysql/main.go