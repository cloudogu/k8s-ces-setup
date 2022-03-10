# Build the manager binary
FROM golang:1.17-alpine3.15 as builder
RUN apk add --no-cache build-base git

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

ENV USER_ID=15000

WORKDIR /

RUN apk add --no-cache bash \
    && set -x \
    && addgroup -g ${USER_ID} -S setup \
    && adduser -u ${USER_ID} -S setup -G setup

# the linter has a problem with the valid colon-syntax
# dockerfile_lint - ignore
USER ${USER_ID}:${USER_ID}

COPY --chown=${USER_ID} resources /
COPY --chown=${USER_ID} --from=builder /workspace/target/k8s-ces-setup /k8s-ces-setup

EXPOSE 8080

ENTRYPOINT ["/startup.sh"]
