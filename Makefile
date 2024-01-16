.PHONY: clean build packing
test:
	@go clean -testcache;
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | tail -n 1 | awk '{ print $3 }'
	@go tool cover -html=coverage.out

external-dep: 
	@go install github.com/cosmtrek/air@latest
	@go install github.com/rubenv/sql-migrate/...@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: lint-fix
## lint-fix: checking code with golangci-lint and fix it
lint-fix:
	@golangci-lint run --fix

.PHONY: generate-swagger
## generate-swagger: generate swagger file
generate-swagger:
	@swag init -g ./cmd/app/app.go