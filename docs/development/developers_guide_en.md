# Developer Guide

This document information to support the development on the `k8s-ces-setup`.

## Local development

First, development files should be created to be used instead of the cluster values:

Dogu-operator-resource:
- place a suitable YAML file (e.g. `dev-dogu-operator.yaml`) under `k8s/dev-resources/`.
- `make serve-local-yaml` returns all resources in the directory
   - Test: [http://localhost:9876/](http://localhost:9876/)
   - a DNS/host alias is useful to communicate from the local K8s cluster to this HTTP server
   - the target needs Python3

`k8s/dev-config/k8s-ces-setup.yaml`:
- `namespace` specifies in which namespace the Cloudogu EcoSystem should be installed
- `dogu_operator_url` specifies the Dogu operator resource
   - e.g. `http://192.168.56.1:9876/dev-dogu-operator.yaml` (see above)

### execution with `go run` or an IDE

- local development at the setup can be started with `STAGE=development go run .`
- execution and debugging in IDEs like IntelliJ is possible
   - but the environment variable `STAGE` should not be forgotten as well

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
