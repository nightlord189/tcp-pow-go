install:
	go mod download

test:
	go test ./...

start-server:
	go run cmd/server/main.go

start-client:
	go run cmd/client/main.go