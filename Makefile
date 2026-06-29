GO_TOOLCHAIN ?= go1.25.11
GO ?= env GOTOOLCHAIN=$(GO_TOOLCHAIN) go
GOFMT ?= $(shell env GOTOOLCHAIN=$(GO_TOOLCHAIN) go env GOROOT 2>/dev/null)/bin/gofmt
TERRAFORM ?= $(shell asdf which terraform 2>/dev/null || command -v terraform)
GOVULNCHECK_VERSION ?= v1.3.0
GOVULNCHECK ?= $(GO) run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
REDOCLY ?= npx --yes @redocly/cli@2.31.5

.PHONY: fmt fmt-check build test test-acceptance bench security openapi-lint docker-build clean review

fmt:
	$(GOFMT) -w cmd internal

fmt-check:
	test -z "$$($(GOFMT) -l cmd internal)"

build:
	$(GO) build -o bin/terraform-provider-terraport ./cmd/terraform-provider-terraport

test:
	TF_ACC_TERRAFORM_PATH="$(TERRAFORM)" $(GO) test ./...

test-acceptance:
	TF_ACC_TERRAFORM_PATH="$(TERRAFORM)" $(GO) test -run 'TestAcc' ./internal/provider

bench:
	$(GO) test -run '^$$' -bench='BenchmarkClient' -benchmem ./internal/bankport

security:
	$(GOVULNCHECK) ./...

openapi-lint:
	$(REDOCLY) lint openapi.yaml

docker-build:
	docker build -t terraport-go-terraform-provider:local .

clean:
	rm -rf bin coverage.out

review:
	$(MAKE) fmt-check
	$(MAKE) build
	$(MAKE) test
	$(MAKE) test-acceptance
	$(MAKE) bench
	$(MAKE) openapi-lint
	$(MAKE) security
