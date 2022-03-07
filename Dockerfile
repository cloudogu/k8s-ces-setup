# Build the manager binary
FROM golang:1.17 as builder

ENV VERSION="0.0.0"

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

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -s -w" -a -o k8s-ces-setup main.go

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
COPY --from=builder /workspace/k8s-ces-setup .

EXPOSE 8080

ENTRYPOINT ["/startup.sh"]
