lint: ## lint code
	golangci-lint run ./...

fix: ## fix backend code
	golangci-lint run --fix ./...
	go fix ./...

go-update: ## update go dependencies
	go get -u $$(go list ./... | grep -v /node_modules)
	go mod tidy

node_modules: package.json pnpm-lock.yaml
	pnpm install
	touch node_modules

public/main.bundle.js: node_modules frontend/src/App.tsx frontend/src/index.tsx webpack.config.js tsconfig.json
	pnpm run build

frontend-install: node_modules ## install frontend dependencies

frontend-build: public/main.bundle.js ## compile frontend code

build: frontend-build ## build cli and api
	go build -o dist/cli cmd/cli/main.go
	go build -o dist/api cmd/api/main.go

test: frontend-test backend-test ## Run all tests

frontend-test: node_modules ## run frontend tests
	pnpm test:frontend

backend-test: ## run backend tests
	go test -count=1 -race ./...

run: frontend-build ## run the service locally
	go run cmd/api/main.go

swagger: ## generate api doc
	swag init -g cmd/api/main.go

help: ## Makefile help
	@grep -E '(^[a-zA-Z_-]+:.*?##.*$$)|(^##)' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'
