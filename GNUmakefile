TEST ?= $$(go list ./...)
GOFMT_FILES ?= $$(find . -name '*.go')
PKG_NAME = akamai

# Local provider install parameters
version = 0.11.0
registry_name = registry.terraform.io
namespace = $(PKG_NAME)
bin_name = terraform-provider-$(PKG_NAME)
build_dir = .build
TF_PLUGIN_DIR ?= ~/.terraform.d/plugins
install_path = $(TF_PLUGIN_DIR)/$(registry_name)/$(namespace)/$(PKG_NAME)/$(version)/$$(go env GOOS)_$$(go env GOARCH)

default: build

.PHONY: install
install: build
	mkdir -p $(install_path)
	cp $(build_dir)/$(bin_name) $(install_path)/$(bin_name)_v$(version)

.PHONY: build
build: fmtcheck
	mkdir -p $(build_dir)
	go build -tags all -o $(build_dir)/$(bin_name)

.PHONY: check
check: errcheck fmtcheck lint vet

.PHONY: test
test: fmtcheck
	go test $(TEST) -v $(TESTARGS) -timeout 10m

.PHONY: testacc
testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 300m

.PHONY: vet
vet:
	@echo "go vet ."
	@go vet $$(go list ./...); if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: fmt
fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY: fmtcheck
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: errcheck
errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

.PHONY: lint
lint: tools.golangci-lint
	@echo "==> Checking source code againse golangci-lint"
	@golangci-lint run ./$(PKG_NAME)

.PHONY: test-compile
test-compile:
	go test -c ./akamai $(TESTARGS)

.PHONY: tools.golangci-lint
tools.golangci-lint:
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint
