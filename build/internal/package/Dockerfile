# syntax=docker/dockerfile:1.0-experimental
ARG TERRAFORM_VERSION="0.13.1"
FROM golang:1.15.1-alpine AS builder
ENV GO111MODULE="on" \
     CGO_ENABLED=0 \
     GOOS="linux" \
     GOARCH="amd64"
ARG PROVIDER_BRANCH_NAME="develop"
WORKDIR $GOPATH/src/github.com/akamai
RUN apk add --update git bash openssh

## Following can only be done inside the lunabuild docker container to copy ssh sockets and keys
COPY .ssh/id_* /root/.ssh/
COPY .ssh/known_hosts /root/.ssh/
RUN chown -R root:root /root/.ssh && chmod 700 /root/.ssh && chmod 600 /root/.ssh/id_rsa && chmod 644 /root/.ssh/id_rsa.pub /root/.ssh/known_hosts

ADD AkamaiCorpRoot-G1.pem /usr/local/share/ca-certificates/AkamaiCorpRoot-G1.pem
RUN update-ca-certificates
RUN git clone ssh://git@git.source.akamai.com:7999/fee/terraform-provider-akamai.git
RUN git clone ssh://git@git.source.akamai.com:7999/fee/akamaiopen-edgegrid-golang.git

RUN cd terraform-provider-akamai && git checkout ${PROVIDER_BRANCH_NAME} && \
    go mod edit -replace github.com/akamai/AkamaiOPEN-edgegrid-golang=../akamaiopen-edgegrid-golang && \
    go install -tags all

FROM hashicorp/terraform:${TERRAFORM_VERSION}
ENV PROVIDER_VERSION="1.0.0"
COPY --from=builder /go/bin/terraform-provider-akamai /root/.terraform.d/plugins/registry.terraform.io/akamai/akamai/${PROVIDER_VERSION}/linux_amd64/terraform-provider-akamai_v${PROVIDER_VERSION}
COPY --from=builder /go/bin/terraform-provider-akamai /root/.terraform.d/plugins/registry.terraform.io/-/akamai/${PROVIDER_VERSION}/linux_amd64/terraform-provider-akamai_v${PROVIDER_VERSION}
