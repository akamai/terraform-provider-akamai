# syntax=docker/dockerfile:1.0-experimental
ARG TERRAFORM_VERSION="0.13.1"
FROM golang:1.15.1-alpine AS builder
ENV GO111MODULE="on" \
     CGO_ENABLED=0 \
     GOOS="linux" \
     GOARCH="amd64"
ARG PROVIDER_BRANCH_NAME="develop"
RUN apk add --update git bash openssh

COPY ./ terraform-provider-akamai
RUN cd terraform-provider-akamai && go install -tags all

FROM hashicorp/terraform:${TERRAFORM_VERSION}
ENV PROVIDER_VERSION="1.0.0"
COPY --from=builder /go/bin/terraform-provider-akamai /root/.terraform.d/plugins/registry.terraform.io/akamai/akamai/${PROVIDER_VERSION}/linux_amd64/terraform-provider-akamai_v${PROVIDER_VERSION}
COPY --from=builder /go/bin/terraform-provider-akamai /root/.terraform.d/plugins/registry.terraform.io/-/akamai/${PROVIDER_VERSION}/linux_amd64/terraform-provider-akamai_v${PROVIDER_VERSION}
