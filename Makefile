# Lets "make" always run test targets
.PHONY: test 
	
build:
	 @go build -o bin/starlight-cargo ./cmd/app/
	
run: build
	@./bin/starlight-cargo

test:
	@go test ./...

