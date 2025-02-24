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