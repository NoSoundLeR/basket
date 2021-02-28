.DEFAULT_GOAL := build

build:
	go build -o ./basket-server ./cmd/server/server.go

run:
	go run ./cmd/server/server.go