TEST?=$$(go list ./... | grep -v /vendor/)
GOFMT_FILES?=$$(find . -name '*.go' | grep -v /vendor/)

default: build

build: fmtcheck
	go install

check: errcheck fmtcheck lint vet

test: fmtcheck
	go test $(TEST) -v $(TESTARGS) -timeout 60s

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

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

.PHONY: build test testacc vet fmt fmtcheck errcheck lint test-compile
