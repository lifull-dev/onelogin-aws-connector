.PHONY: deps install build cross-build test coverage-html clean

NAME = onelogin-aws-connector

deps:
	go mod download

install: deps test
	go install -tags includeClientToken

build: deps test
	mkdir -p build
	go build -o build/$(NAME) -tags includeClientToken

cross-build: deps test
	GOOS=linux GOARCH=amd64 go build -tags includeClientToken -o dist/linux-amd64/$(NAME)
	GOOS=darwin GOARCH=amd64 go build -tags includeClientToken -o dist/darwin-amd64/$(NAME)
	GOOS=windows GOARCH=amd64 go build -tags includeClientToken -o dist/windows-aml64/$(NAME)

test:
	go test ./... -cover

coverage-html:
	mkdir .coverage
	go test ./... -cover -coverprofile=.coverage/coverage.out
	go tool cover -html=.coverage/coverage.out

clean:
	-rm -rf build
	-rm -rf dist
	-rm -rf .coverage
	-rm -rf vendor
