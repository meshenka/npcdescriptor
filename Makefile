lint:
	golangci-lint run ./...

build:
	go build -o cli cmd/cli/main.go 

test:
	go test -count=1 -race ./...
