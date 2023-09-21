build:
	@go build -o go-jwt-auth

run: build
	@./go-jwt-auth

test:
	@go test -v ./...