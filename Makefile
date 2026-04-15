lint:
	golangci-lint run ./...

fix:
	golangci-lint run --fix ./...
	go fix ./...

frontend-install:
	pnpm install

frontend-build:
	pnpm run build

build: frontend-build
	go build -o cli cmd/cli/main.go
	go build -o api cmd/api/main.go

test:
	go test -count=1 -race ./...

run: frontend-build
	go run cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go
