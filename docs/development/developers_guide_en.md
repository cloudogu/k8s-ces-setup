# Developer Guide

This document information to support the development on the `k8s-ces-setup`.

## Local development

First, development files should be created to be used instead of the cluster values:

Dogu-operator-resource:
- place a suitable YAML file (e.g. `dev-dogu-operator.yaml`) under `k8s/dev-resources/`.
- `make serve-local-yaml` returns all resources in the directory
   <!-- markdown-link-check-disable-next-line -->
   - Test: http://localhost:9876/
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

## Debugging

It is possible to interact with a cluster-deployed setup:

```bash
# Check setup state
curl --request GET --url http://192.168.56.2:30080/api/v1/health
{"status":"healthy","version":"0.0.0"}

# Create a namespace according the Setup configuration map
curl -I --request POST --url http://192.168.56.2:30080/api/v1/setup
```

## Restore pre-setup state

Sometimes it is necessary to turn the time back to the beginning, e.g. to check installation routines. This can be done with the following commands (pay attention to your **current namespace**):

```bash
# delete the resources directly created by the setup
make k8s-delete
# deletes target namespace and all namespaced resources in it (pods, deployments, secrets, etc.)
kubectl delete ns your-namespace
# deletes CRD so that it can be initially imported with the dogu operator
kubectl delete crd dogus.k8s.cloudogu.com
# deletes clusterroles/bindings from setup installations
kubectl delete clusterroles k8s-dogu-operator-metrics-reader ingress-nginx
kubectl delete clusterrolebindings ingress-nginx
# manually delete resources that may still be deployed incorrectly
...
```

## Cleanup of the setup.

When new Kubernetes resources are created in development, they may need to be taken into account by the cleanup task.
To do this, edit the configmap `k8s-ces-setup-cleanup-script` script in `k8s-ces-setup.yaml`.