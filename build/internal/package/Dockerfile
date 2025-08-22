FROM golang:1.23.11-alpine3.21

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
    go install github.com/axw/gocov/gocov@v1.1.0 && \
    go install github.com/AlekSi/gocov-xml@v1.1.0 && \
    go install github.com/jstemmer/go-junit-report@v1.0.0 && \
    mkdir -p /root/.ssh

ADD build/internal/package/AkamaiCorpRoot-G1.pem /usr/local/share/ca-certificates/AkamaiCorpRoot-G1.pem
RUN update-ca-certificates
