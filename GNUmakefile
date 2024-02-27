TEST ?= $$(go list ./...)
PKG_NAME = akamai

# Local provider install parameters
version = 0.11.0
registry_name = registry.terraform.io
namespace = $(PKG_NAME)
bin_name = terraform-provider-$(PKG_NAME)
build_dir = .build
TF_PLUGIN_DIR ?= ~/.terraform.d/plugins
install_path = $(TF_PLUGIN_DIR)/$(registry_name)/$(namespace)/$(PKG_NAME)/$(version)/$$(go env GOOS)_$$(go env GOARCH)

BIN      = $(CURDIR)/bin
GOCMD = go
GOTEST = $(GOCMD) test
GOBUILD = $(GOCMD) build
M = $(shell echo ">")
TFLINT = $(BIN)/tflint
$(BIN)/tflint: $(BIN) ; $(info $(M) Installing tflint...)
	@export TFLINT_INSTALL_PATH=$(BIN); \
	curl -sSfL https://raw.githubusercontent.com/terraform-linters/tflint/master/install_linux.sh  | bash

$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) Installing $(PACKAGE)...)
	@tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) $(GOCMD) get $(PACKAGE) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

GOIMPORTS = $(BIN)/goimports
$(BIN)/goimports: PACKAGE=golang.org/x/tools/cmd/goimports

GOLANGCI_LINT_VERSION = v1.55.2
GOLANGCILINT = $(BIN)/golangci-lint
$(BIN)/golangci-lint: ; $(info $(M) Installing golangci-lint...) @
	$Q curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN) $(GOLANGCI_LINT_VERSION)

# Targets
default: build

.PHONY: install
install: build
	mkdir -p $(install_path)
	cp $(build_dir)/$(bin_name) $(install_path)/$(bin_name)_v$(version)

.PHONY: build
build:
	mkdir -p $(build_dir)
	go build -o $(build_dir)/$(bin_name)

.PHONY: tidy
tidy:
	@go mod tidy
	@cd tools && go mod tidy

.PHONY: test
test:
	go test $(TEST) -v $(TESTARGS) -timeout 40m 2>&1

.PHONY: testacc
testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 300m

.PHONY: fmt
fmt:  | $(GOIMPORTS); $(info $(M) Running goimports...) @ ## Run goimports on all source files
	$Q $(GOIMPORTS) -w .

.PHONY: fmt-check
fmt-check: | $(GOIMPORTS); $(info $(M) Running format and imports check...) @ ## Run goimports on all source files
	$(eval OUTPUT = $(shell $(GOIMPORTS) -l .))
	@if [ "$(OUTPUT)" != "" ]; then\
		echo "Found following files with incorrect format and/or imports:";\
		echo "$(OUTPUT)";\
		false;\
	fi

.PHONY: terraform-fmtcheck
terraform-fmtcheck:
	terraform fmt -recursive -check

.PHONY: terraform-fmt
terraform-fmt:
	terraform fmt -recursive

.PHONY: lint
lint: | $(GOLANGCILINT) ; $(info $(M) Running golangci-lint...) @
	go mod tidy
	$Q $(BIN)/golangci-lint run

.PHONY: terraform-lint
terraform-lint: | $(TFLINT) ; $(info $(M) Checking source code against tflint...) @ ## Run tflint on all HCL files in the project
	@find ./examples -type f -name "*.tf" | xargs -I % dirname % | sort -u | xargs -I @ sh -c "echo @ && $(TFLINT) --filter @"

.PHONY: test-compile
test-compile:
	go test -c ./akamai $(TESTARGS)

.PHONY: tools
tools: $(TOOLS)

.PHONY: clean-tools
clean-tools:
	@rm -rf $(TOOLS_BIN_DIR)

.PHONY: init
init: tools tools.terraform

.PHONY: tools.terraform
tools.terraform:
	@sh -c "'$(CURDIR)/scripts/install_terraform.sh'"
