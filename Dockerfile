# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY app app

# Copy git and build files
COPY .git .git
COPY Makefile Makefile
COPY build build

# Build
RUN go mod vendor
RUN make compile-generic

## Production image
FROM alpine:3.15.0
LABEL maintainer="hello@cloudogu.com" \
      NAME="k8s-ces-setup" \
      VERSION="0.0.0"

ENV USER=setup

WORKDIR /

RUN apk add --no-cache bash \
    && set -x \
    && addgroup -S ${USER} \
    && adduser -S ${USER} -G ${USER}

# the linter has a problem with the valid colon-syntax
# dockerfile_lint - ignore
USER ${USER}:${USER}

COPY resources /
COPY --from=builder /workspace/target/k8s-ces-setup .

EXPOSE 8080

ENTRYPOINT ["/startup.sh"]
