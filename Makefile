DOCKER=docker

BUILD_DIR=build
CMD_DIR=cmd
CMDS=$(patsubst $(CMD_DIR)/%,%,$(wildcard $(CMD_DIR)/*))

.PHONY: fmt check test

all: fmt check test

fmt:
	gofmt -s -w -l .
	@echo 'goimports' && goimports -w -local gobe $(shell find . -type f -name '*.go' -not -path "./internal/*")
	gci write -s standard -s default --skip-generated .
	@files=$$(git diff --diff-filter=AM --name-only HEAD^ | grep '.go'); \
	if [ -n "$$files" ]; then \
	  echo 'golines' $$files; \
		golines --ignore-generated -m 80 --reformat-tags --shorten-comments -w $$files; \
	fi
	go mod tidy

check:
	revive -exclude pkg/... -formatter friendly -config .revive.toml  ./...
	go vet -all ./...
	golangci-lint run
	misspell -error */**
	@echo 'staticcheck' && staticcheck $(shell go list ./... | grep -v internal)

test:
	go test ./...

