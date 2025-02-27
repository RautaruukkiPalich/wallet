.PHONY: run

run:
	go run ./cmd/wallet/main.go

lint:
	golangci-lint run ./...

pprof:
	go tool pprof [binary] http://127.0.0.1:8081/debug/pprof/profile

swagger:
	swag init -g ./cmd/wallet/main.go
	swag fmt


## DATABASE
MIGRATIONS_DB_URI := "postgres://postgres:postgres@localhost:35432/wallet?sslmode=disable"
MIGRATIONS_PATH := ./migrations

newmigrate:
	migrate create -ext sql -dir ./migrations -seq $(NAME)

migrateup:
	migrate -database $(MIGRATIONS_DB_URI) -path $(MIGRATIONS_PATH) up

migratedown:
	migrate -database $(MIGRATIONS_DB_URI) -path $(MIGRATIONS_PATH) down

## BOMBING
bomb-get-amount:
	bombardier -c 50 -d 30s http://localhost:8080/api/v1/wallets/00000000-0000-0000-0000-000000000001

bomb-withdraw:
	bombardier -c 2 -d 30s -m POST  http://localhost:8080/api/v1/wallets/ -H 'accept: application/json' -H 'Content-Type: application/json' -f './example-withdraw.json'

bomb-deposit:
	bombardier -c 2 -d 30s -m POST  http://localhost:8080/api/v1/wallets/ -H 'accept: application/json' -H 'Content-Type: application/json' -f './example-deposit.json'

bomb-new-wallet:
	bombardier -c 50 -d 30s -m POST  http://localhost:8080/api/v1/wallets/create -H 'accept: application/json'
