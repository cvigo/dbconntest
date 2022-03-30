FROM ubuntu:latest

MAINTAINER Titan Team <iaasdbs.group@bbva.com>

## Install dependencies
RUN apt-get update \
    && DEBIAN_FRONTEND="noninteractive" apt-get install -y \
        make \
        uuid-runtime \
        wget \
        curl \
        gcc \
        git \
        libxml2-dev \
    && rm -rf /var/lib/apt/lists/*

# Install Golang
RUN wget -c https://dl.google.com/go/go1.15.linux-amd64.tar.gz -O - | tar -xz -C /opt \
    && mv /opt/go /usr/local

## Golang environment variables
ENV GOROOT=/usr/local/go \
    GOCACHE=/cache/go-build \
    GOLANGCI_LINT_CACHE=/cache/go \
    GOPATH=/go \
    CGO_ENABLED="1" \
    DISABLE_WARN_OUTSIDE_CONTAINER=1 \
    GOPRIVATE="*.bbva.com,*.igrupobbva" \
    GOPROXY="https://proxy.golang.org,direct" \
    GO111MODULE=on \
    GOMODCACHE="$GOPATH/pkg/mod"

## Include Golang binaries path
ENV PATH=${PATH}:${GOROOT}/bin:${GOPATH}/bin

## Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.38.0

## Install libraries required by go_ibm_db module
RUN GO111MODULE=off go get -d github.com/ibmdb/go_ibm_db \
    && cd ${GOPATH}/src/github.com/ibmdb/go_ibm_db/installer \
    && go run setup.go \
    && cd -

## IBM DB environment variables
ENV DB2HOME="${GOPATH}/src/github.com/ibmdb/clidriver"
ENV LD_LIBRARY_PATH="$DB2HOME/lib" \
    CGO_CFLAGS="-I$DB2HOME/include" \
    CGO_LDFLAGS="-L$DB2HOME/lib"

## Certs needed to verify the identity of the Git server
COPY  ssl/certs /etc/ssl/certs/

RUN git config --global url."git@scm.live.es.nextgen.igrupobbva:".insteadOf "https://scm.live.es.nextgen.igrupobbva/"

## Create a workspace
WORKDIR /workspace

## Default command
CMD ["go", "test", "./..."]
