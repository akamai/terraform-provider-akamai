# syntax=docker/dockerfile:1.0-experimental
ARG TERRAFORM_VERSION
FROM golang:1.15.1-alpine AS builder
ENV GO111MODULE="on" \
     CGO_ENABLED=0 \
     GOOS="linux" \
     GOARCH="amd64"
ARG PROVIDER_BRANCH_NAME
WORKDIR $GOPATH/src/github.com/akamai
RUN apk add --update git bash openssh

RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan git.source.akamai.com >> ~/.ssh/known_hosts
RUN --mount=type=ssh git clone ssh://git@git.source.akamai.com:7999/fee/terraform-provider-akamai.git
RUN --mount=type=ssh git clone ssh://git@git.source.akamai.com:7999/fee/akamaiopen-edgegrid-golang.git
RUN cd terraform-provider-akamai \
    && git checkout ${PROVIDER_BRANCH_NAME} \
    && go mod edit -replace github.com/akamai/AkamaiOPEN-edgegrid-golang=../akamaiopen-edgegrid-golang \
    && go install -tags all

FROM hashicorp/terraform:${TERRAFORM_VERSION}
#WORKDIR tmp/terraform
ENV PROVIDER_VERSION="1.0.0"
COPY --from=builder /go/bin/terraform-provider-akamai /root/.terraform.d/plugins/registry.terraform.io/akamai/akamai/${PROVIDER_VERSION}/linux_amd64/terraform-provider-akamai_v${PROVIDER_VERSION}
COPY --from=builder /go/bin/terraform-provider-akamai /root/.terraform.d/plugins/registry.terraform.io/-/akamai/${PROVIDER_VERSION}/linux_amd64/terraform-provider-akamai_v${PROVIDER_VERSION}