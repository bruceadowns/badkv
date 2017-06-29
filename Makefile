all: check test build

check: goimports govet

goimports:
	@goimports -d .

govet:
	@go tool vet -all .

test:
	@go test -v -cover

build:
	@go build

run:
	@go run --race main.go

clean:
	@-rm -v ./badkv
