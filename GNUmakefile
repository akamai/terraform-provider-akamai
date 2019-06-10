TEST?=$$(go list ./... | grep -v /vendor/)
GOFMT_FILES?=$$(find . -name '*.go' | grep -v /vendor/)

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

lint:
	@sh -c "'$(CURDIR)/scripts/golint.sh'"

test-compile:
	go test -c ./akamai $(TESTARGS)

dep:
	@which dep > /dev/null; if [ $$? -ne 0 ]; then \
		echo "==> Installing dep..."; \
		go get -u github.com/golang/dep/cmd/dep; \
	fi

dep-install: dep
	@echo "==> Installing vendor dependencies..."
	@dep ensure -vendor-only

dep-update: dep
	@echo "==> Updating vendor dependencies..."
	@dep ensure -update

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

.PHONY: build build-docker test test-docker testacc vet fmt fmtcheck errcheck test-compile website website-test
