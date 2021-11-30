# syntax=docker/dockerfile:1.0-experimental
ARG TERRAFORM_VERSION="1.0.4"
FROM hashicorp/terraform:${TERRAFORM_VERSION}
ENV PROVIDER_VERSION="1.0.0" \
    GO111MODULE="on" \
    CGO_ENABLED=0 \
    GOOS="linux" \
    GOARCH="amd64" \
    PATH=$PATH:/root/go/bin

ARG SSH_PRV_KEY
ARG SSH_PUB_KEY
ARG SSH_KNOWN_HOSTS
WORKDIR $GOPATH/src/github.com/akamai

RUN apk add --update git bash sudo openssh go gcc musl-dev openssl-dev ca-certificates unzip curl && \
    go get -u github.com/jstemmer/go-junit-report && \
    go get -u github.com/axw/gocov/... && \
    go get -u github.com/AlekSi/gocov-xml && \
    mkdir -p /root/.ssh && \
    curl -s https://raw.githubusercontent.com/terraform-linters/tflint/master/install_linux.sh | bash

ADD build/internal/package/AkamaiCorpRoot-G1.pem /usr/local/share/ca-certificates/AkamaiCorpRoot-G1.pem
RUN update-ca-certificates
