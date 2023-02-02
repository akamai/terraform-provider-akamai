# syntax=docker/dockerfile:1.0-experimental
ARG TERRAFORM_VERSION="1.2.5"
FROM alpine:3.16
ENV PROVIDER_VERSION="1.0.0" \
    CGO_ENABLED=0 \
    GOOS="linux" \
    GOARCH="amd64" \
    PATH=$PATH:/usr/local/go/bin:/root/go/bin \
    TFLINT_VERSION="v0.39.3"

ARG GOLANGCI_LINT_VERSION="v1.50.1"
ARG SSH_PRV_KEY
ARG SSH_PUB_KEY
ARG SSH_KNOWN_HOSTS
WORKDIR $GOPATH/src/github.com/akamai

RUN apk add --update git bash sudo openssh gcc go musl-dev openssl-dev ca-certificates unzip curl terraform && \
    go install github.com/axw/gocov/gocov@latest && \
    go install github.com/AlekSi/gocov-xml@latest && \
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION} && \
    go install github.com/jstemmer/go-junit-report@latest && \
    mkdir -p /root/.ssh && \
    curl -s https://raw.githubusercontent.com/terraform-linters/tflint/master/install_linux.sh | bash

ADD build/internal/package/AkamaiCorpRoot-G1.pem /usr/local/share/ca-certificates/AkamaiCorpRoot-G1.pem
RUN update-ca-certificates
