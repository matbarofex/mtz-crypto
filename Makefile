SVC_NAME := $(shell grep 'const ServiceName' pkg/version.go | cut -f2 -d '"')
APP_VSN := $(shell grep 'const Version' pkg/version.go | cut -f2 -d '"')
GIT_COMMIT := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date --iso-8601=seconds)
IMAGE_NAME := mtzio/${SVC_NAME}

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

.PHONY: docker-build
docker-build:
	docker build \
		--build-arg SVC_NAME="$(SVC_NAME)" \
		--build-arg APP_VSN="$(APP_VSN)" \
		--build-arg GIT_COMMIT="$(GIT_COMMIT)" \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		-t $(IMAGE_NAME):$(APP_VSN)-$(GIT_COMMIT) \
		-t $(IMAGE_NAME):latest \
		.

.PHONY: docker-push
docker-push:
	docker push $(IMAGE_NAME):$(APP_VSN)-$(GIT_COMMIT)
	docker push $(IMAGE_NAME):latest
