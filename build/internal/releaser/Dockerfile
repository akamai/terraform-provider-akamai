FROM golang:1.25.7-alpine3.23

WORKDIR /workspace

ENV DEBIAN_FRONTEND=noninteractive
ENV PATH=$PATH:/usr/local/go/bin:/root/go/bin

ARG SSH_PRV_KEY
ARG SSH_PUB_KEY
ARG SSH_KNOWN_HOSTS
ARG TERRAFORM_VERSION="1.13.5"

COPY AkamaiCorpRoot-G1.pem /usr/local/share/ca-certificates/AkamaiCorpRoot-G1.pem
COPY .edgerc /root/.edgerc
COPY clone_repos.bash /usr/local/bin/clone_repos.bash
COPY goreleaser_build.bash /usr/local/bin/goreleaser_build.bash
COPY smoke_tests.bash /usr/local/bin/smoke_tests.bash

RUN apk add --update git bash openssh gcc musl-dev openssl-dev ca-certificates make curl && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    mv terraform /usr/bin/terraform \
    && go install github.com/goreleaser/goreleaser/v2@v2.5.1 \
    && mkdir -p /root/.terraform.d/plugins/registry.terraform.io/akamai/akamai/999.0.0/linux_amd64 /root/.ssh

