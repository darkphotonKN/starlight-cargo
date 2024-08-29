# Lets "make" always run test targets
.PHONY: test 
	
build:
	 @go build -o bin/fs ./cmd/app/
	
run: build
	@./bin/fs

test:
	@go test ./...

