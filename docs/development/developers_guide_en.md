# Developer Guide

This document information to support the development on the `k8s-ces-setup`.

## Local development

Local development on the setup can be started with `go run .`.

## Makefile targets

The command `make help` prints all available targets and their descriptions on the command line.

For the Makefiles to work with respect to the cluster, the root path of the development environment must be entered in the
Makefiles under the environment variable `K8S_CLUSTER_ROOT`.

## Debugging

It is possible to interact with a cluster-deployed setup:

```bash
# Check setup state
curl --request GET --url http://192.168.56.2:30080/api/v1/health
{"status":"healthy","version":"0.0.0"}

# Create a namespace according the Setup configuration map
curl -I --request POST --url http://192.168.56.2:30080/api/v1/setup
```
