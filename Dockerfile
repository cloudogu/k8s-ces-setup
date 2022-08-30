# Build the manager binary
FROM golang:1.18.1 as builder

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

# Copy .git files as the build process builds the current commit id into the binary via ldflags
COPY .git .git

# Copy build files
COPY build build
COPY Makefile Makefile

# Build
RUN go mod vendor
RUN make compile-generic

## Production image
FROM gcr.io/distroless/static:nonroot
LABEL maintainer="hello@cloudogu.com" \
      NAME="k8s-ces-setup" \
      VERSION="0.7.1"

WORKDIR /

# the linter has a problem with the valid colon-syntax
# dockerfile_lint - ignore
USER 15000:15000

COPY --chown=15000:15000 --from=builder /workspace/target/k8s-ces-setup /k8s-ces-setup

EXPOSE 8080

ENTRYPOINT ["/k8s-ces-setup"]
