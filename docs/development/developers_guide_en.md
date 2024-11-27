# Developer Guide

This document information to support the development on the `k8s-ces-setup`.

## Local development

First, development files should be created to be used instead of the cluster values:

* `k8s/dev-resources/k8s-ces-setup.yaml`: [setup-config](../operations/configuration_guide_en.md)
* `k8s/dev-resources/setup.json`: [custom-setup-config](../operations/custom_setup_configuration_en.md)


### Installing the ces setup in the local cluster

In order for the ces-setup to be executed and tested in the local cluster, a few things must be taken into account.
Firstly, all existing dogus, components, etc. should be removed from the system. To do this
the command `make k8s-clean` can be used.
So that the Ces setup can then be installed, a small change must first be made to the
k8s/helm/values.yaml file beforehand.
The following part must be commented in, otherwise the setup cannot be carried out:
```
  # k8s-longhorn:
  # version: latest
  # helmRepositoryNamespace: k8s
  # deployNamespace: longhorn-system
```
The Ces setup can then be installed with `make helm-apply`. It is then carried out automatically.

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

Sometimes it is necessary to turn the time back to the beginning, e.g. to check installation routines.
This can be done with the make target `k8s-clean` (pay attention to the **current namespace**):

```bash
# delete all dogus & components and all the resources directly created by the setup
make k8s-clean

# manually delete any resources that may still have been incorrectly created
...
```

## Cleanup of the setup.

When new Kubernetes resources are created in development, they may need to be taken into account by the cleanup task.
To do this, edit the configmap `k8s-ces-setup-cleanup-script` script in `k8s-ces-setup.yaml`.