FROM ubuntu:focal

WORKDIR /workspace

ENV DEBIAN_FRONTEND=noninteractive
ENV PATH=$PATH:/usr/local/go/bin:/root/go/bin

ARG SSH_PRV_KEY
ARG SSH_PUB_KEY
ARG SSH_KNOWN_HOSTS

COPY AkamaiCorpRoot-G1.pem /usr/local/share/ca-certificates/AkamaiCorpRoot-G1.pem
COPY .edgerc /root/.edgerc
COPY clone_repos.bash /usr/local/bin/clone_repos.bash
COPY goreleaser_build.bash /usr/local/bin/goreleaser_build.bash
COPY smoke_tests.bash /usr/local/bin/smoke_tests.bash

RUN apt update && apt install -y curl git gcc ca-certificates openssh-client gnupg \
    && echo "deb [arch=amd64] https://apt.releases.hashicorp.com focal main" >> /etc/apt/sources.list \
    && curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add - \
    && apt update && apt install -y terraform \
    && update-ca-certificates \
    && curl -o go1.18.linux-amd64.tar.gz https://dl.google.com/go/go1.18.linux-amd64.tar.gz \
    && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz \
    && go install github.com/goreleaser/goreleaser@latest \
    && mkdir -p /root/.terraform.d/plugins/registry.terraform.io/akamai/akamai/10.0.0/linux_amd64 /root/.ssh

