.DEFAULT_GOAL := build

build:
	go build ./cmd/server/server.go

run: build
	./server.exe