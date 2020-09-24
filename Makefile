ROOT := github.com/mehdy/twixter

GO = $(GO_VARS) go

GO_VARS ?= GOOS=linux GOARCH=amd64

.PHONY: mocks

mocks:
	$(GO) generate ./...

