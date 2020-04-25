TEST?=$$(go list ./... | grep -v /vendor/)
GOFMT_FILES?=$$(find . -name '*.go' | grep -v /vendor/)
PKG_NAME=akamai
WEBSITE_REPO=github.com/hashicorp/terraform-website

default: build

build: fmtcheck
	go install

check: errcheck fmtcheck lint vet

test: fmtcheck
	go test $(TEST) -v $(TESTARGS) -timeout 60s

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 300m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v /vendor/); if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

#lint:
#	@echo "==> Checking source code against linters..."
#	@golangci-lint run ./$(PKG_NAME)

lint: tools.golangci-lint
	@echo "==> Checking source code againse golangci-lint"
	@golangci-lint run ./$(PKG_NAME)

tools:
	@echo "==> installing required tooling..."
#	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
#	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	#GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint

test-compile:
	go test -c ./akamai $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build build-docker test test-docker testacc vet fmt fmtcheck errcheck lint test-compile website website-test

.PHONY: tools.golangci-lint

tools.golangci-lint:
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint
