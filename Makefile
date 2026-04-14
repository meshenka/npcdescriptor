lint:
	golangci-lint run ./...

fix:
	golangci-lint run --fix ./...
	go fix ./...

build:
	go build -o cli cmd/cli/main.go
	go build -o api cmd/api/main.go

test:
	go test -count=1 -race ./...

run:
	go run cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go
