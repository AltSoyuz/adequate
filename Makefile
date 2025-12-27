VERSION := 0.0.1
BUILDINFO_PKG := $(shell go list -f '{{.ImportPath}}' ./lib/buildinfo)

vet:
	go vet ./...

fmt: 
	find . -type f -name '*.go' -not -path './vendor/*' | xargs gofmt -l -w -s

check-all: check-gitleaks fmt vet golangci-lint 

ui-install-dependencies:
	cd ui && npm ci

check-ui: ui-install-dependencies
	cd ui && npm run lint

build-ui: ui-install-dependencies
	cd ui && npm run build

test-short:
	go test -short ./cmd/... ./internal/... ./lib/... 

test-race:
	go test -race ./cmd/... ./internal/... ./lib/...

test-full:
	go test -coverprofile=coverage.txt -covermode=atomic ./cmd/... ./internal/... ./lib/...

integration-test: clean build
	go test ./apptest/... 

build: clean
	go build -ldflags "-X $(BUILDINFO_PKG).Version=$(VERSION)" -o ./bin/app ./cmd/main.go

vendor-update:
	go get -u ./...
	go mod tidy 
	go mod vendor

golangci-lint: install-golangci-lint
	golangci-lint run

install-golangci-lint:
	which golangci-lint || \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(shell go env GOPATH)/bin v2.5.0

remove-golangci-lint:
	rm -rf `which golangci-lint`

clean:
	rm -rf bin/*

sqlc:
	docker run --rm -v ./:/src -w /src sqlc/sqlc generate

check-gitleaks:
	docker run --rm \
		-v ./:/src \
		-w /src \
		ghcr.io/gitleaks/gitleaks:latest detect \
		--verbose \
		--redact=100 \
		--report-path=gitleaks-report.json \
		--report-format=json \
		--exit-code=1
	