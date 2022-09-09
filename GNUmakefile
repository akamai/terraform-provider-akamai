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

# Tools versions
golangci-lint-version = v1.41.1
tflint-version        = v0.39.3 # Newer versions contain rules that examples are not compliant with

default: build

.PHONY: install
install: build
	mkdir -p $(install_path)
	cp $(build_dir)/$(bin_name) $(install_path)/$(bin_name)_v$(version)

.PHONY: build
build:
	mkdir -p $(build_dir)
	go build -tags all -o $(build_dir)/$(bin_name)

.PHONY: check
check: errcheck fmtcheck lint vet

.PHONY: test
test:
	go test $(TEST) -v $(TESTARGS) -timeout 20m 2>&1 

.PHONY: testacc
testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 300m

.PHONY: vet
vet:
	@echo "==> Checking source code against vet"
	# Appsec package excluded until https://track.akamai.com/jira/browse/SECKSD-12824 is done
	@go vet $$(go list ./... | grep -v appsec); if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: fmt
fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY: terraform-fmt
terraform-fmt:
	terraform fmt -recursive -check

.PHONY: fmtcheck
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: errcheck
errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

.PHONY: lint
lint:
	@echo "==> Checking source code against golangci-lint"
	@$$(go env GOPATH)/bin/golangci-lint run

.PHONY: terraform-lint
terraform-lint:
	@echo "==> Checking source code against tflint"
	@find ./examples -type f -name "*.tf" | xargs -I % dirname % | sort -u | xargs -I @ sh -c "echo @ && tflint @"

.PHONY: test-compile
test-compile:
	go test -c ./akamai $(TESTARGS)

.PHONY: tools.golangci-lint
tools.golangci-lint:
	@echo Installing golangci-lint
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(golangci-lint-version)

.PHONY: tools.tflint
tools.tflint:
	@echo Installing tf-lint
	@export TFLINT_VERSION=$(tflint-version) && curl -s https://raw.githubusercontent.com/terraform-linters/tflint/master/install_linux.sh | bash

.PHONY: init
init: tools.golangci-lint tools.tflint

.PHONY: dummy-edgerc
dummy-edgerc:
	@sh -c "'$(CURDIR)/scripts/dummyedgerc.sh'"

.PHONY: tools.terraform
tools.terraform:
	@sh -c "'$(CURDIR)/scripts/install_terraform.sh'"
