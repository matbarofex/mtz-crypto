.PHONY: build
build:
	go build -o mtz-crypto-service ./cmd/...

.PHONY: build-race
build-race:
	go build -race -o mtz-crypto-service ./cmd/...

.PHONY: generate-mocks
generate-mocks:
	mockery --all

.PHONY: test
test: build
	go test -covermode=atomic -race -v -count=1 -coverprofile=coverage.out ./pkg/...
	go tool cover -func coverage.out | grep total

.PHONY: lint
lint:
	golangci-lint run

.PHONY: load-test
load-test:
	h2load --h1 -c 50 -t 4 -D 10 --warm-up-time=2s -i test-urls.txt -B 'http://localhost:8000'
