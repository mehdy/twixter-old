ROOT := github.com/mehdy/twixter

GO = $(GO_VARS) go

GO_VARS ?= GOOS=linux GOARCH=amd64
.PHONY: deps tests mocks

deps:
	$(GO) mod tidy

mocks:
	$(GO) generate ./...


tests:
	$(GO) test -v -race -cover ./...
