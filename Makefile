PATH := ${PWD}/bin:${PATH}
export PATH

# linters
bin/reviewdog:
	curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s -- -b ./bin v0.9.13

bin/golangci-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ./bin v1.19.0

#  Commands
#----------------------------------------------------------------

.PHONY: lint
lint: ./bin/reviewdog ./bin/golangci-lint
ifdef CI
	reviewdog -reporter=github-pr-review
else
	reviewdog -diff="git diff master"
endif

.PHONY: test
test:
	go test -v ./...

.PHONY: cover
cover:
	go test -v -coverprofile coverage.txt -covermode atomic ./...
