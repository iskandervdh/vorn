GOBIN ?= $$(go env GOPATH)/bin
PLATFORMS = darwin/amd64 darwin/arm64 linux/386 linux/amd64 linux/arm linux/arm64 windows/386 windows/amd64 windows/arm64

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml

.PHONY: run-all-examples
run-all-examples:
	set -e
	# Run all examples
	for filename in ./examples/*.vorn; do \
		./vorn "$$filename"; \
	done

.PHONY: build
build:
	for platform in ${PLATFORMS}; do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} go build -o build/vorn-$${platform%/*}-$${platform#*/}; \
	done
