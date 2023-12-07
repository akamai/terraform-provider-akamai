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

# Developer tools
TOOLS_MOD_FILE := $(CURDIR)/tools/go.mod
TOOLS_BIN_DIR := $(CURDIR)/tools/bin
TOOL_PKGS := $(shell go list -f '{{join .Imports " "}}' tools/tools.go)
TOOLS := $(foreach TOOL,$(notdir $(TOOL_PKGS)),$(TOOLS_BIN_DIR)/$(TOOL))
$(foreach TOOL,$(TOOLS),$(eval $(notdir $(TOOL)) := $(TOOL))) # Allows to use e.g. $(golangci-lint) instead of $(TOOLS_BIN_DIR)/golangci-lint

$(TOOLS_BIN_DIR):
	@mkdir -p $(TOOLS_BIN_DIR)

$(TOOLS_MOD_FILE): tidy

$(TOOLS): $(TOOLS_MOD_FILE) | $(TOOLS_BIN_DIR)
	$(eval TOOL := $(filter %/$(@F),$(TOOL_PKGS)))
	$(eval TOOL_VERSION := $(shell grep -m 1 $(shell echo $(TOOL) | cut -d/ -f 1-3) $(TOOLS_MOD_FILE) | cut -d' ' -f2))
	@echo "Installing $(TOOL)@$(TOOL_VERSION)"
	@GOBIN=$(TOOLS_BIN_DIR) go install -modfile=$(TOOLS_MOD_FILE) $(TOOL)

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
	go test $(TEST) -v $(TESTARGS) -timeout 30m 2>&1

.PHONY: testacc
testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 300m

.PHONY: fmt
fmt: $(goimports)
	$(goimports) -w .

.PHONY: terraform-fmtcheck
terraform-fmtcheck:
	terraform fmt -recursive -check

.PHONY: terraform-fmt
terraform-fmt:
	terraform fmt -recursive

.PHONY: lint
lint: $(golangci-lint)
	$(golangci-lint) run

.PHONY: terraform-lint
terraform-lint: $(tflint)
	@find ./examples -type f -name "*.tf" | xargs -I % dirname % | sort -u | xargs -I @ sh -c "echo @ && $(tflint) @"

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
