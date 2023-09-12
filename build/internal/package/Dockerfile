FROM golang:1.20-alpine3.18

ENV PROVIDER_VERSION="1.0.0" \
    CGO_ENABLED=0 \
    GOOS="linux" \
    GOARCH="amd64" \
    PATH=$PATH:/usr/local/go/bin:/root/go/bin

ARG SSH_PRV_KEY
ARG SSH_PUB_KEY
ARG SSH_KNOWN_HOSTS
WORKDIR $GOPATH/src/github.com/akamai

RUN apk add --update git bash sudo openssh gcc go musl-dev openssl-dev ca-certificates unzip curl make && \
    go install github.com/axw/gocov/gocov@latest && \
    go install github.com/AlekSi/gocov-xml@latest && \
    go install github.com/jstemmer/go-junit-report@latest && \
    mkdir -p /root/.ssh

ADD build/internal/package/AkamaiCorpRoot-G1.pem /usr/local/share/ca-certificates/AkamaiCorpRoot-G1.pem
RUN update-ca-certificates
