lint:
	golangci-lint run ./...

fix:
	golangci-lint run --fix ./...
	go fix ./...

node_modules: package.json pnpm-lock.yaml
	pnpm install
	touch node_modules

public/main.bundle.js: node_modules frontend/src/App.tsx frontend/src/index.tsx webpack.config.js tsconfig.json
	pnpm run build

frontend-install: node_modules

frontend-build: public/main.bundle.js

build: frontend-build
	go build -o cli cmd/cli/main.go
	go build -o api cmd/api/main.go

test: frontend-test backend-test

frontend-test: node_modules
	pnpm test:frontend

backend-test:
	go test -count=1 -race ./...

run: frontend-build
	go run cmd/api/main.go


swagger:
	swag init -g cmd/api/main.go
