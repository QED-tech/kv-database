run-server: build server
test: gen unit

install:
	go mod tidy

build:
	rm bin/client bin/server 2>/dev/null || true
	go build -race -o bin/server cmd/server/main.go && chmod +x bin/server
	go build -race -o bin/client cmd/client/main.go && chmod +x bin/client

server:
	./bin/server

client:
	./bin/client

gen:
	go generate ./...

unit:
	go test --race ./internal/...